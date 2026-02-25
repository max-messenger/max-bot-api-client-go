package maxbot

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/stretchr/testify/require"
)

func Test_uploads_UploadMediaFromReader_whenUploadVideoOrAudioType(t *testing.T) {
	var server *httptest.Server

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/uploads" {
			uploadType := r.FormValue("type")

			switch uploadType {
			case string(schemes.VIDEO):
				_, _ = fmt.Fprintf(w, `{"token": "new_video_token", "url": "%s/mock-upload-media"}`, server.URL)
			case string(schemes.AUDIO):
				_, _ = fmt.Fprintf(w, `{"token": "new_audio_token", "url": "%s/mock-upload-media"}`, server.URL)
			case string(schemes.FILE):
				_, _ = fmt.Fprintf(w, `{"url": "%s/mock-upload-file"}`, server.URL)
			}
			return
		}

		if r.URL.Path == "/mock-upload-media" {
			_, _ = fmt.Fprint(w, "<retval>1</retval>")
			return
		}

		if r.URL.Path == "/mock-upload-file" {
			_, _ = fmt.Fprint(w, `{"file_id": 12345, "token": "new_file_token"}`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	tests := []struct {
		name       string
		uploadType schemes.UploadType
		want       *schemes.UploadedInfo
	}{
		{
			name:       "video type ignores xml and returns token",
			uploadType: schemes.VIDEO,
			want:       &schemes.UploadedInfo{Token: "new_video_token"},
		},
		{
			name:       "audio type ignores xml and returns token",
			uploadType: schemes.AUDIO,
			want:       &schemes.UploadedInfo{Token: "new_audio_token"},
		},
		{
			name:       "file type parses json correctly",
			uploadType: schemes.FILE,
			want:       &schemes.UploadedInfo{FileID: 12345, Token: "new_file_token"},
		},
	}

	u, _ := url.Parse(server.URL)
	cl := newClient("bot_token", Version, u, server.Client())
	upl := newUploads(cl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := upl.UploadMediaFromReader(context.Background(), tt.uploadType, strings.NewReader("fake file content"))

			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}
