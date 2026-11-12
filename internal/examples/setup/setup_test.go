package setup

import (
	"os"
	"strings"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		need    []string
		want    Env
		wantErr string // подстрока сообщения; "" => ошибки нет
	}{
		{
			name: "отрицательный chat_id сохраняется со знаком и шириной",
			env:  map[string]string{KeyToken: "t", KeyChat: "-70801090403050"},
			need: []string{KeyToken, KeyChat},
			want: Env{Token: "t", Chat: -70801090403050},
		},
		{
			name: "chat_id без минуса разбирается",
			env:  map[string]string{KeyToken: "t", KeyChat: "70801090403050"},
			need: []string{KeyToken, KeyChat},
			want: Env{Token: "t", Chat: 70801090403050},
		},
		{
			name: "chat_id вне int32 не портится",
			env:  map[string]string{KeyToken: "t", KeyChat: "4294967296"},
			need: []string{KeyToken, KeyChat},
			want: Env{Token: "t", Chat: 4294967296},
		},
		{
			name:    "пустой BOT_TOKEN даёт инструкцию, а не панику",
			env:     map[string]string{KeyChat: "-1"},
			need:    []string{KeyToken, KeyChat},
			wantErr: KeyToken + " не задан",
		},
		{
			name:    "нечисловой chat_id называет ключ и формат",
			env:     map[string]string{KeyToken: "t", KeyChat: "abc"},
			need:    []string{KeyToken, KeyChat},
			wantErr: KeyChat,
		},
		{
			name:    "отсутствующий chat_id называет действие",
			env:     map[string]string{KeyToken: "t"},
			need:    []string{KeyToken, KeyChat},
			wantErr: "добавьте в неё бота",
		},
		{
			name:    "отсутствующий TEST_MEMBER_ID называет роль B",
			env:     map[string]string{KeyToken: "t"},
			need:    []string{KeyToken, KeyMember},
			wantErr: KeyMember,
		},
		{
			name:    "неизвестный ключ окружения",
			env:     map[string]string{KeyToken: "t"},
			need:    []string{"NOPE"},
			wantErr: "неизвестный ключ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got, err := LoadEnv(tt.need...)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("неожиданная ошибка: %v", err)
				}
				if got != tt.want {
					t.Fatalf("env = %+v, хотим %+v", got, tt.want)
				}
				return
			}
			if err == nil {
				t.Fatalf("ожидалась ошибка с %q, получили env %+v", tt.wantErr, got)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("ошибка %q не содержит %q", err, tt.wantErr)
			}
		})
	}
}

func TestProofString(t *testing.T) {
	for _, p := range []Proof{L1, L2, L3} {
		if s := p.String(); s == "" || strings.HasPrefix(s, "L?") {
			t.Errorf("Proof(%d) не имеет подписи: %q", p, s)
		}
	}
}

// TestFixture проверяет, что фикстура находится и при CWD=корень модуля
// (`go run`), и при CWD=каталог пакета (`go test`).
func TestFixture(t *testing.T) {
	got, err := Fixture("test_data/test_file.txt")
	if err != nil {
		t.Fatalf("Fixture: %v", err)
	}
	if _, err := os.Stat(got); err != nil {
		t.Fatalf("Fixture вернул несуществующий путь %q: %v", got, err)
	}
	if _, err := Fixture("test_data/нет-такого-файла"); err == nil {
		t.Fatal("для отсутствующей фикстуры ожидалась ошибка")
	}
}
