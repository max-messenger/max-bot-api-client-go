package maxbot

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type uploads struct {
	client *client
}

func newUploads(client *client) *uploads {
	return &uploads{client: client}
}

// UploadMediaFromFile uploads the file to the Max server.
func (a *uploads) UploadMediaFromFile(ctx context.Context, uploadType schemes.UploadType, filename string) (*schemes.UploadedInfo, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer fh.Close()

	return a.UploadMediaFromReaderWithName(ctx, uploadType, fh, filename)
}

// UploadMediaFromUrl uploads the file from a remote server to the Max server.
// urlStr is the URL of the file to download (e.g. "https://example.com/file.pdf").
func (a *uploads) UploadMediaFromUrl(ctx context.Context, uploadType schemes.UploadType, urlStr string) (*schemes.UploadedInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := a.client.do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch URL: %w", err)
	}
	defer resp.Body.Close()

	return a.UploadMediaFromReaderWithName(ctx, uploadType, resp.Body, a.attachmentName(resp))
}

func (a *uploads) UploadMediaFromReader(ctx context.Context, uploadType schemes.UploadType, reader io.Reader) (*schemes.UploadedInfo, error) {
	result := new(schemes.UploadedInfo)

	return result, a.uploadMediaFromReader(ctx, uploadType, reader, "", result)
}

func (a *uploads) UploadMediaFromReaderWithName(ctx context.Context, uploadType schemes.UploadType, reader io.Reader, name string) (*schemes.UploadedInfo, error) {
	result := new(schemes.UploadedInfo)

	return result, a.uploadMediaFromReader(ctx, uploadType, reader, name, result)
}

// UploadPhotoFromFile uploads photos to the Max server.
func (a *uploads) UploadPhotoFromFile(ctx context.Context, fileName string) (*schemes.PhotoTokens, error) {
	fh, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer fh.Close()
	result := new(schemes.PhotoTokens)

	return result, a.uploadMediaFromReader(ctx, schemes.PHOTO, fh, fileName, result)
}

// UploadPhotoFromBase64String uploads photos to the Max server.
func (a *uploads) UploadPhotoFromBase64String(ctx context.Context, code string) (*schemes.PhotoTokens, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(code))
	result := new(schemes.PhotoTokens)

	return result, a.uploadMediaFromReader(ctx, schemes.PHOTO, decoder, "", result)
}

// UploadPhotoFromUrl uploads the photo from a remote server to the Max server.
// urlStr is the URL of the image to download.
func (a *uploads) UploadPhotoFromUrl(ctx context.Context, urlStr string) (*schemes.PhotoTokens, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	respFile, err := a.client.do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch URL: %w", err)
	}
	defer respFile.Body.Close()
	result := new(schemes.PhotoTokens)
	name := a.attachmentName(respFile)

	return result, a.uploadMediaFromReader(ctx, schemes.PHOTO, respFile.Body, name, result)
}

// UploadPhotoFromReader uploads the photo from a reader.
func (a *uploads) UploadPhotoFromReader(ctx context.Context, reader io.Reader) (*schemes.PhotoTokens, error) {
	result := new(schemes.PhotoTokens)

	return result, a.uploadMediaFromReader(ctx, schemes.PHOTO, reader, "", result)
}

func (a *uploads) UploadPhotoFromReaderWithName(ctx context.Context, reader io.Reader, name string) (*schemes.PhotoTokens, error) {
	result := new(schemes.PhotoTokens)

	return result, a.uploadMediaFromReader(ctx, schemes.PHOTO, reader, name, result)
}

func (a *uploads) getUploadURL(ctx context.Context, uploadType schemes.UploadType) (*schemes.UploadEndpoint, error) {
	result := new(schemes.UploadEndpoint)
	values := url.Values{}
	values.Set(paramType, string(uploadType))
	body, err := a.client.request(ctx, http.MethodPost, pathUpload, values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err := json.NewDecoder(body).Decode(result); err != nil {
		return result, fmt.Errorf("decode upload endpoint: %w", err)
	}
	return result, nil
}

func (a *uploads) uploadMediaFromReader(
	ctx context.Context,
	uploadType schemes.UploadType,
	reader io.Reader,
	fileName string,
	result interface{},
) error {
	endpoint, err := a.getUploadURL(ctx, uploadType)
	if err != nil {
		return err
	}
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	if fileName == "" {
		fileName = "file"
	} else {
		fileName = filepath.Base(fileName)
	}

	fileWriter, err := bodyWriter.CreateFormFile("data", fileName)
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	if _, err = io.Copy(fileWriter, reader); err != nil {
		return fmt.Errorf("copy file data: %w", err)
	}

	contentType := bodyWriter.FormDataContentType()
	if err := bodyWriter.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.Url, bodyBuf)
	if err != nil {
		return fmt.Errorf("create upload request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := a.client.do(req)
	if err != nil {
		return fmt.Errorf("upload: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &schemes.Error{}
		if decodeErr := json.NewDecoder(resp.Body).Decode(apiErr); decodeErr == nil {
			return &APIError{
				Code:    resp.StatusCode,
				Message: apiErr.Code,
				Details: apiErr.Message,
			}
		}
		return fmt.Errorf("upload: HTTP %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode upload response: %w", err)
	}

	return nil
}

func (*uploads) attachmentName(r *http.Response) string {
	disposition := r.Header["Content-Disposition"]
	if len(disposition) != 0 {
		_, params, err := mime.ParseMediaType(disposition[0])
		if err == nil && params["filename"] != "" {
			return params["filename"]
		}
	}

	return ""
}
