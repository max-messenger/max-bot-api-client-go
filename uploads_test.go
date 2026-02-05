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
				fmt.Fprint(w, `{"token": "new_video_token", "url": "`+server.URL+`/mock-upload-video-or-audio"}`)
			case string(schemes.AUDIO):
				fmt.Fprint(w, `{"token": "new_audio_token", "url": "`+server.URL+`/mock-upload-video-or-audio"}`)
			case string(schemes.FILE):
				fmt.Fprint(w, `{"url": "`+server.URL+`/mock-upload-file"}`)
			}
		}
		if r.URL.Path == "/mock-upload-video-or-audio" {
			fmt.Fprint(w, "<retval>1</retval>")
		}
		if r.URL.Path == "/mock-upload-file" {
			fmt.Fprint(w, `{"file_id": 12345, "token": "new_file_token"}`)
		}
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	cl := newClient("bot_token", version, u, server.Client())
	upl := newUploads(cl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := upl.UploadMediaFromReader(context.Background(), tt.uploadType, strings.NewReader("content"))
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}
