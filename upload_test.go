package maxbot

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type mockHttpClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if m.doFunc != nil {
		return m.doFunc(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil
}

func TestUpload_Upload(t *testing.T) {
	tests := []struct {
		name           string
		uploadType     model.UploadType
		fileContent    string
		fileName       string
		fileSize       int64
		setupMock      func(*mockHttpClient)
		expectedToken  string
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:        "successful upload of image",
			uploadType:  model.UploadImage,
			fileContent: "fake image data",
			fileName:    "test.jpg",
			fileSize:    16,
			setupMock: func(m *mockHttpClient) {
				// Mock getUploadURL response
				m.doFunc = func(req *http.Request) (*http.Response, error) {
					if strings.Contains(req.URL.Path, pathUpload) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewReader([]byte(`{
								"url": "https://upload.example.com/upload"
							}`))),
						}, nil
					}
					// Mock actual upload response
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(bytes.NewReader([]byte(`{
							"token": "uploaded_token_123"
						}`))),
					}, nil
				}
			},
			expectedToken: "uploaded_token_123",
			expectedError: false,
		},
		{
			name:        "successful upload of audio without token in response",
			uploadType:  model.UploadAudio,
			fileContent: "fake audio data",
			fileName:    "test.mp3",
			fileSize:    16,
			setupMock: func(m *mockHttpClient) {
				m.doFunc = func(req *http.Request) (*http.Response, error) {
					if strings.Contains(req.URL.Path, pathUpload) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewReader([]byte(`{
								"url": "https://upload.example.com/upload",
								"token": "audio_token_456"
							}`))),
						}, nil
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte{})),
					}, nil
				}
			},
			expectedToken: "audio_token_456",
			expectedError: false,
		},
		{
			name:        "successful upload of video without token in response",
			uploadType:  model.UploadVideo,
			fileContent: "fake video data",
			fileName:    "test.mp4",
			fileSize:    16,
			setupMock: func(m *mockHttpClient) {
				m.doFunc = func(req *http.Request) (*http.Response, error) {
					if strings.Contains(req.URL.Path, pathUpload) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewReader([]byte(`{
								"url": "https://upload.example.com/upload",
								"token": "video_token_789"
							}`))),
						}, nil
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte{})),
					}, nil
				}
			},
			expectedToken: "video_token_789",
			expectedError: false,
		},
		{
			name:        "upload with empty filename uses default",
			uploadType:  model.UploadFile,
			fileContent: "file content",
			fileName:    "",
			fileSize:    12,
			setupMock: func(m *mockHttpClient) {
				m.doFunc = func(req *http.Request) (*http.Response, error) {
					if strings.Contains(req.URL.Path, pathUpload) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewReader([]byte(`{
								"url": "https://upload.example.com/upload"
							}`))),
						}, nil
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(bytes.NewReader([]byte(`{
							"token": "default_name_token"
						}`))),
					}, nil
				}
			},
			expectedToken: "default_name_token",
			expectedError: false,
		},
		{
			name:        "error getting upload URL",
			uploadType:  model.UploadImage,
			fileContent: "test data",
			fileName:    "test.txt",
			fileSize:    9,
			setupMock: func(m *mockHttpClient) {
				m.doFunc = func(req *http.Request) (*http.Response, error) {
					if strings.Contains(req.URL.Path, pathUpload) {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body: io.NopCloser(bytes.NewReader([]byte(`{
								"error": "upload_url_failed"
							}`))),
						}, nil
					}
					return nil, nil
				}
			},
			expectedError:  true,
			expectedErrMsg: "getUploadURL err",
		},
		{
			name:        "upload returns non-OK status",
			uploadType:  model.UploadImage,
			fileContent: "test data",
			fileName:    "test.txt",
			fileSize:    9,
			setupMock: func(m *mockHttpClient) {
				m.doFunc = func(req *http.Request) (*http.Response, error) {
					if strings.Contains(req.URL.Path, pathUpload) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewReader([]byte(`{
								"url": "https://upload.example.com/upload"
							}`))),
						}, nil
					}
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body: io.NopCloser(bytes.NewReader([]byte(`{
							"error": "upload_failed"
						}`))),
					}, nil
				}
			},
			expectedError: true,
		},
		{
			name:        "invalid JSON response for image",
			uploadType:  model.UploadImage,
			fileContent: "test data",
			fileName:    "test.txt",
			fileSize:    9,
			setupMock: func(m *mockHttpClient) {
				m.doFunc = func(req *http.Request) (*http.Response, error) {
					if strings.Contains(req.URL.Path, pathUpload) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewReader([]byte(`{
								"url": "https://upload.example.com/upload"
							}`))),
						}, nil
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`invalid json`))),
					}, nil
				}
			},
			expectedError: true,
		},
		{
			name:        "upload with very large file",
			uploadType:  model.UploadFile,
			fileContent: strings.Repeat("a", 1024*1024), // 1MB
			fileName:    "large.txt",
			fileSize:    1024 * 1024,
			setupMock: func(m *mockHttpClient) {
				m.doFunc = func(req *http.Request) (*http.Response, error) {
					if strings.Contains(req.URL.Path, pathUpload) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewReader([]byte(`{
								"url": "https://upload.example.com/upload"
							}`))),
						}, nil
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(bytes.NewReader([]byte(`{
							"token": "large_file_token"
						}`))),
					}, nil
				}
			},
			expectedToken: "large_file_token",
			expectedError: false,
		},
		{
			name:        "upload with Chinese filename",
			uploadType:  model.UploadFile,
			fileContent: "test content",
			fileName:    "文件名.txt",
			fileSize:    12,
			setupMock: func(m *mockHttpClient) {
				m.doFunc = func(req *http.Request) (*http.Response, error) {
					if strings.Contains(req.URL.Path, pathUpload) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewReader([]byte(`{
								"url": "https://upload.example.com/upload"
							}`))),
						}, nil
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(bytes.NewReader([]byte(`{
							"token": "chinese_filename_token"
						}`))),
					}, nil
				}
			},
			expectedToken: "chinese_filename_token",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHttpClient{}
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			cli := &client{
				token:       "test_token",
				httpClient:  mockClient,
				pollPause:   time.Second,
				pollTimeout: 30 * time.Second,
			}
			cli.baseURL.Scheme = defaultScheme
			cli.baseURL.Host = "api.example.com"

			upload := newUpload(cli)

			ctx := context.Background()
			reader := strings.NewReader(tt.fileContent)

			token, err := upload.Upload(ctx, tt.uploadType, reader, tt.fileName, tt.fileSize)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.expectedErrMsg != "" && !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token != tt.expectedToken {
					t.Errorf("Expected token '%s', got '%s'", tt.expectedToken, token)
				}
			}
		})
	}
}

func TestUpload_Upload_ContextCancellation(t *testing.T) {
	mockClient := &mockHttpClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			// Simulate slow upload by waiting for context cancellation
			<-req.Context().Done()
			return nil, req.Context().Err()
		},
	}

	cli := &client{
		token:      "test_token",
		httpClient: mockClient,
	}
	cli.baseURL.Scheme = defaultScheme
	cli.baseURL.Host = "api.example.com"

	upload := newUpload(cli)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	reader := strings.NewReader("test data")
	_, err := upload.Upload(ctx, model.UploadImage, reader, "test.txt", 9)

	if err == nil {
		t.Errorf("Expected context cancellation error, got nil")
	}
}

func TestUpload_Upload_EmptyReader(t *testing.T) {
	mockClient := &mockHttpClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.Path, pathUpload) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewReader([]byte(`{
						"url": "https://upload.example.com/upload"
					}`))),
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader([]byte(`{
					"token": "empty_file_token"
				}`))),
			}, nil
		},
	}

	cli := &client{
		token:      "test_token",
		httpClient: mockClient,
	}
	cli.baseURL.Scheme = defaultScheme
	cli.baseURL.Host = "api.example.com"

	upload := newUpload(cli)

	ctx := context.Background()
	reader := strings.NewReader("")

	token, err := upload.Upload(ctx, model.UploadFile, reader, "empty.txt", 0)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if token != "empty_file_token" {
		t.Errorf("Expected token 'empty_file_token', got '%s'", token)
	}
}

// Test multipartFileName helper function
func TestMultipartFileName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty filename", "", defaultFileName},
		{"simple filename", "test.jpg", "test.jpg"},
		{"path with directories", "path/to/file.pdf", "file.pdf"},
		{"filename with spaces", "my file name.doc", "my file name.doc"},
		{"unicode filename", "测试文件.pdf", "测试文件.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := multipartFileName(tt.input)
			if result != tt.expected {
				t.Errorf("multipartFileName(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Test multipartEnvelope helper function
func TestMultipartEnvelope(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		fileSize int64
	}{
		{"simple file", "test.txt", 100},
		{"empty name", "", 50},
		{"large file", "large.bin", 1024 * 1024},
		{"filename with spaces", "my file.pdf", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentType, contentLength, boundary, err := multipartEnvelope(tt.fileName, tt.fileSize)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if contentType == "" {
				t.Errorf("ContentType should not be empty")
			}

			if contentLength <= tt.fileSize {
				t.Errorf("ContentLength (%d) should be greater than fileSize (%d) due to multipart overhead", contentLength, tt.fileSize)
			}

			if boundary == "" {
				t.Errorf("Boundary should not be empty")
			}
		})
	}
}
