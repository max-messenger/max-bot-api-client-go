package maxbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

// generateTestHash генерирует валидный hash для тестовых данных
func generateTestHash(botToken string, params map[string]string) string {
	// Создаем копию параметров и удаляем hash
	values := url.Values{}
	for k, v := range params {
		if k != "hash" {
			values.Set(k, v)
		}
	}

	// Сортируем параметры
	var sortedParams []string
	for key := range values {
		value := values.Get(key)
		sortedParams = append(sortedParams, fmt.Sprintf("%s=%s", key, value))
	}
	sort.Strings(sortedParams)

	dataCheckString := strings.Join(sortedParams, "\n")

	mac1 := hmac.New(sha256.New, []byte("WebAppData"))
	mac1.Write([]byte(botToken))
	secretKey := mac1.Sum(nil)

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(dataCheckString))

	return hex.EncodeToString(mac.Sum(nil))
}

// createValidInitData создает валидные init данные для тестирования
func createValidInitData(botToken string, userData any) string {
	userJSON, _ := json.Marshal(userData)

	params := map[string]string{
		"user":          string(userJSON),
		"chat_instance": "test_chat_123",
		"chat_type":     "private",
		"start_param":   "test_start",
		"auth_date":     fmt.Sprintf("%d", time.Now().Unix()),
	}

	hash := generateTestHash(botToken, params)
	params["hash"] = hash

	// Формируем query string
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	return values.Encode()
}

func TestValidateInitData(t *testing.T) {
	validBotToken := "test_bot_token_123"
	validUser := model.UserApp{
		ID:        123456,
		FirstName: "Test",
		LastName:  "User",
		Username:  "testuser",
	}

	tests := []struct {
		name       string
		initData   string
		botToken   string
		wantErr    bool
		wantErrMsg string
		wantUser   model.UserApp
	}{
		{
			name:     "Valid init data",
			initData: createValidInitData(validBotToken, validUser),
			botToken: validBotToken,
			wantErr:  false,
			wantUser: validUser,
		},
		{
			name:       "Empty init data",
			initData:   "",
			botToken:   validBotToken,
			wantErr:    true,
			wantErrMsg: "initData cannot be empty",
		},
		{
			name:       "Empty bot token",
			initData:   "some_data",
			botToken:   "",
			wantErr:    true,
			wantErrMsg: "botToken cannot be empty",
		},
		{
			name:       "Missing hash parameter",
			initData:   "user=%7B%22id%22%3A123%7D&auth_date=123456789",
			botToken:   validBotToken,
			wantErr:    true,
			wantErrMsg: "hash parameter is missing",
		},
		{
			name: "Invalid hash verification",
			initData: func() string {
				params := map[string]string{
					"user":      `{"id":123}`,
					"auth_date": "123456789",
					"hash":      "invalid_hash_value",
				}
				values := url.Values{}
				for k, v := range params {
					values.Set(k, v)
				}
				return values.Encode()
			}(),
			botToken:   validBotToken,
			wantErr:    true,
			wantErrMsg: "hash verification failed",
		},
		{
			name: "Invalid JSON in user parameter",
			initData: func() string {
				params := map[string]string{
					"user":      "invalid_json",
					"auth_date": "123456789",
				}
				hash := generateTestHash(validBotToken, params)
				params["hash"] = hash
				values := url.Values{}
				for k, v := range params {
					values.Set(k, v)
				}
				return values.Encode()
			}(),
			botToken:   validBotToken,
			wantErr:    true,
			wantErrMsg: "json decode err",
		},
		{
			name: "URL encoded init data",
			initData: func() string {
				rawData := createValidInitData(validBotToken, validUser)
				return url.QueryEscape(rawData)
			}(),
			botToken: validBotToken,
			wantErr:  false,
			wantUser: validUser,
		},
		{
			name: "User with missing fields",
			initData: func() string {
				userWithMissingFields := struct {
					ID int64 `json:"id"`
				}{ID: 789}
				return createValidInitData(validBotToken, userWithMissingFields)
			}(),
			botToken: validBotToken,
			wantErr:  false,
			wantUser: model.UserApp{ID: 789},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := ValidateInitData(tt.initData, tt.botToken)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateInitData() expected error but got none")
					return
				}
				if tt.wantErrMsg != "" && !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("ValidateInitData() error = %v, want error containing %v", err, tt.wantErrMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateInitData() unexpected error: %v", err)
				return
			}

			if gotUser != tt.wantUser {
				t.Errorf("ValidateInitData() user = %+v, want %+v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestValidateInitDataWithWebAppPlatform(t *testing.T) {
	botToken := "test_bot_token"
	user := model.UserApp{ID: 123, FirstName: "Test"}

	t.Run("WebAppPlatform parameter should be ignored in hash", func(t *testing.T) {
		params := map[string]string{
			"user":             mustMarshalJSON(user),
			"auth_date":        fmt.Sprintf("%d", time.Now().Unix()),
			"web_app_platform": "ios",
			"web_app_version":  "1.0",
			"chat_instance":    "test123",
		}

		hash := generateTestHash(botToken, params)
		params["hash"] = hash

		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		initData := values.Encode()

		gotUser, err := ValidateInitData(initData, botToken)
		if err != nil {
			t.Errorf("ValidateInitData() failed with web_app_platform: %v", err)
		}

		if gotUser != user {
			t.Errorf("ValidateInitData() user = %+v, want %+v", gotUser, user)
		}
	})
}

func TestValidateInitDataWithSpecialCharacters(t *testing.T) {
	botToken := "test_token_!@#$%"
	user := model.UserApp{
		ID:        123,
		FirstName: "Test & Special",
		LastName:  "User=Value",
		Username:  "test@user.com",
	}

	t.Run("Special characters in user data", func(t *testing.T) {
		initData := createValidInitData(botToken, user)

		gotUser, err := ValidateInitData(initData, botToken)
		if err != nil {
			t.Errorf("ValidateInitData() failed with special characters: %v", err)
		}

		if gotUser != user {
			t.Errorf("ValidateInitData() user = %+v, want %+v", gotUser, user)
		}
	})
}

func TestValidateInitDataTampering(t *testing.T) {
	botToken := "test_bot_token"
	originalUser := model.UserApp{ID: 123, FirstName: "Original"}

	t.Run("Tampered user data should fail verification", func(t *testing.T) {
		initData := createValidInitData(botToken, originalUser)

		// Изменяем данные пользователя
		modifiedUser := model.UserApp{ID: 999, FirstName: "Hacked"}
		modifiedUserJSON, _ := json.Marshal(modifiedUser)

		values, _ := url.ParseQuery(initData)
		values.Set("user", string(modifiedUserJSON))
		tamperedData := values.Encode()

		_, err := ValidateInitData(tamperedData, botToken)
		if err == nil {
			t.Error("ValidateInitData() should fail with tampered data")
		}

		if !strings.Contains(err.Error(), "hash verification failed") {
			t.Errorf("Expected hash verification failed, got: %v", err)
		}
	})

	t.Run("Tampered auth_date should fail verification", func(t *testing.T) {
		initData := createValidInitData(botToken, originalUser)

		values, _ := url.ParseQuery(initData)
		values.Set("auth_date", "9999999999")
		tamperedData := values.Encode()

		_, err := ValidateInitData(tamperedData, botToken)
		if err == nil {
			t.Error("ValidateInitData() should fail with tampered auth_date")
		}
	})
}

func TestValidateInitDataMultipleHashValues(t *testing.T) {
	botToken := "test_bot_token"
	user := model.UserApp{ID: 123}

	t.Run("Multiple hash parameters - should use first one", func(t *testing.T) {
		params := map[string]string{
			"user":      mustMarshalJSON(user),
			"auth_date": fmt.Sprintf("%d", time.Now().Unix()),
		}

		validHash := generateTestHash(botToken, params)

		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		values.Add("hash", validHash)
		values.Add("hash", "extra_hash")
		initData := values.Encode()

		gotUser, err := ValidateInitData(initData, botToken)
		if err != nil {
			t.Errorf("ValidateInitData() failed: %v", err)
		}

		if gotUser != user {
			t.Errorf("ValidateInitData() user = %+v, want %+v", gotUser, user)
		}
	})
}

func BenchmarkValidateInitData(b *testing.B) {
	botToken := "test_bot_token"
	user := model.UserApp{ID: 123456, FirstName: "Test", LastName: "User"}
	initData := createValidInitData(botToken, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ValidateInitData(initData, botToken)
	}
}

// Helper function
func mustMarshalJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
