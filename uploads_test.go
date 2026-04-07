package maxbot

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/stretchr/testify/require"
)

func Test_uploads_UploadMediaFromReaderWithName_whenUploadVideoOrAudioType(t *testing.T) {
	tests := []struct {
		name       string
		uploadType schemes.UploadType
		want       *schemes.UploadedInfo
	}{
		{
			name:       "video type",
			uploadType: schemes.VIDEO,
			want:       &schemes.UploadedInfo{Token: "new_video_token"},
		},
		{
			name:       "audio type",
			uploadType: schemes.AUDIO,
			want:       &schemes.UploadedInfo{Token: "new_audio_token"},
		},
		{
			name:       "file type",
			uploadType: schemes.FILE,
			want:       &schemes.UploadedInfo{FileID: 12345, Token: "new_file_token"},
		},
	}

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/uploads" {
			switch r.URL.Query().Get("type") {
			case string(schemes.VIDEO):
				_, _ = fmt.Fprint(w, `{"token": "new_video_token", "url": "`+server.URL+`/mock-upload-video-or-audio"}`)
			case string(schemes.AUDIO):
				_, _ = fmt.Fprint(w, `{"token": "new_audio_token", "url": "`+server.URL+`/mock-upload-video-or-audio"}`)
			case string(schemes.FILE):
				_, _ = fmt.Fprint(w, `{"url": "`+server.URL+`/mock-upload-file"}`)
			}
		}
		if r.URL.Path == "/mock-upload-video-or-audio" {
			_, _ = fmt.Fprint(w, "<retval>1</retval>")
		}
		if r.URL.Path == "/mock-upload-file" {
			_, _ = fmt.Fprint(w, `{"file_id": 12345, "token": "new_file_token"}`)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	cl := newClient("bot_token", Version, u, server.Client())
	upl := newUploads(cl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := upl.UploadMediaFromReader(context.Background(), tt.uploadType, strings.NewReader("content"))
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func Test_uploads_UploadMediaFromFile(t *testing.T) {
	const fileContent = "file content"

	tempDir := t.TempDir()
	fileName := filepath.Join(tempDir, "payload.txt")
	require.NoError(t, os.WriteFile(fileName, []byte(fileContent), 0o600))

	_, contentLength, _, err := multipartEnvelope(filepath.Base(fileName), int64(len(fileContent)))
	require.NoError(t, err)

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/uploads":
			_, _ = fmt.Fprint(w, `{"url": "`+server.URL+`/mock-upload-file"}`)
		case "/mock-upload-file":
			require.True(t, strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary="))
			require.NotEmpty(t, r.Header.Get("Content-Type"))
			require.Equal(t, contentLength, r.ContentLength)

			part, err := r.MultipartReader()
			require.NoError(t, err)

			filePart, err := part.NextPart()
			require.NoError(t, err)
			require.Equal(t, "data", filePart.FormName())
			require.Equal(t, filepath.Base(fileName), filePart.FileName())

			body, err := io.ReadAll(filePart)
			require.NoError(t, err)
			require.Equal(t, fileContent, string(body))

			_, err = part.NextPart()
			require.ErrorIs(t, err, io.EOF)

			_, _ = fmt.Fprint(w, `{"file_id": 12345, "token": "new_file_token"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	cl := newClient("bot_token", Version, u, server.Client())
	upl := newUploads(cl)

	result, err := upl.UploadMediaFromFile(context.Background(), schemes.FILE, fileName)
	require.NoError(t, err)
	require.Equal(t, &schemes.UploadedInfo{FileID: 12345, Token: "new_file_token"}, result)
}

func Test_uploads_UploadPhotoFromFile(t *testing.T) {
	const fileContent = "photo content"

	tempDir := t.TempDir()
	fileName := filepath.Join(tempDir, "photo.jpg")
	require.NoError(t, os.WriteFile(fileName, []byte(fileContent), 0o600))

	_, contentLength, _, err := multipartEnvelope(filepath.Base(fileName), int64(len(fileContent)))
	require.NoError(t, err)

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/uploads":
			require.Equal(t, string(schemes.PHOTO), r.URL.Query().Get("type"))
			_, _ = fmt.Fprint(w, `{"url": "`+server.URL+`/mock-upload-photo"}`)
		case "/mock-upload-photo":
			require.True(t, strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary="))
			require.NotEmpty(t, r.Header.Get("Content-Type"))
			require.Equal(t, contentLength, r.ContentLength)
			_, _ = fmt.Fprint(w, `{"photos":{"default":{"token":"photo_token"}}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	cl := newClient("bot_token", Version, u, server.Client())
	upl := newUploads(cl)

	result, err := upl.UploadPhotoFromFile(context.Background(), fileName)
	require.NoError(t, err)
	require.Equal(t, &schemes.PhotoTokens{
		Photos: map[string]schemes.PhotoToken{
			"default": {Token: "photo_token"},
		},
	}, result)
}

func Test_uploads_UploadMediaFromFile_matchesBufferedMultipartHeaders(t *testing.T) {
	const fileContent = "same content"

	tempDir := t.TempDir()
	fileName := filepath.Join(tempDir, "payload.txt")
	require.NoError(t, os.WriteFile(fileName, []byte(fileContent), 0o600))

	type uploadRequest struct {
		header        http.Header
		contentLength int64
	}

	requests := make([]uploadRequest, 0, 2)

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/uploads":
			_, _ = fmt.Fprint(w, `{"url": "`+server.URL+`/mock-upload-file"}`)
		case "/mock-upload-file":
			requests = append(requests, uploadRequest{
				header:        r.Header.Clone(),
				contentLength: r.ContentLength,
			})

			reader, err := r.MultipartReader()
			require.NoError(t, err)

			part, err := reader.NextPart()
			require.NoError(t, err)

			body, err := io.ReadAll(part)
			require.NoError(t, err)
			require.Equal(t, fileContent, string(body))

			_, _ = fmt.Fprint(w, `{"file_id": 12345, "token": "new_file_token"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	cl := newClient("bot_token", Version, u, server.Client())
	upl := newUploads(cl)

	_, err := upl.UploadMediaFromReaderWithName(context.Background(), schemes.FILE, strings.NewReader(fileContent), fileName)
	require.NoError(t, err)

	_, err = upl.UploadMediaFromFile(context.Background(), schemes.FILE, fileName)
	require.NoError(t, err)

	require.Len(t, requests, 2)

	bufferedType, bufferedParams, err := mime.ParseMediaType(requests[0].header.Get("Content-Type"))
	require.NoError(t, err)
	streamingType, streamingParams, err := mime.ParseMediaType(requests[1].header.Get("Content-Type"))
	require.NoError(t, err)

	require.Equal(t, bufferedType, streamingType)
	require.Equal(t, "multipart/form-data", streamingType)
	require.NotEmpty(t, bufferedParams["boundary"])
	require.NotEmpty(t, streamingParams["boundary"])
	require.Equal(t, requests[0].contentLength, requests[1].contentLength)
}
