package maxbot

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"

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
	defer a.client.closer("uploadMediaFromFile file", fh)

	info, err := fh.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}
	result := new(schemes.UploadedInfo)

	return result, a.uploadMediaFromReaderWithSize(ctx, uploadType, fh, filename, info.Size(), result)
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
	defer a.client.closer("uploadMediaFromUrl body", resp.Body)

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
	defer a.client.closer("uploadPhotoFromFile file", fh)

	info, err := fh.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}
	result := new(schemes.PhotoTokens)

	return result, a.uploadMediaFromReaderWithSize(ctx, schemes.PHOTO, fh, fileName, info.Size(), result)
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

	defer a.client.closer("uploadPhotoFromUrl body", respFile.Body)
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
	defer a.client.closer("getUploadURL body", body)

	return result, jsoniter.NewDecoder(body).Decode(result)
}

func (a *uploads) uploadMediaFromReader(
	ctx context.Context,
	uploadType schemes.UploadType,
	reader io.Reader,
	fileName string,
	result any,
) error {
	endpoint, err := a.getUploadURL(ctx, uploadType)
	if err != nil {
		return err
	}
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileName = multipartFileName(fileName)

	fileWriter, err := bodyWriter.CreateFormFile("data", fileName)
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	if _, err = io.Copy(fileWriter, reader); err != nil {
		return fmt.Errorf("copy file data: %w", err)
	}

	contentType := bodyWriter.FormDataContentType()
	if err = bodyWriter.Close(); err != nil {
		return err
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
	defer a.client.closer("uploadMediaFromReader body", resp.Body)

	return a.decodeUploadResponse(resp, uploadType, endpoint.Token, result)
}

func (a *uploads) uploadMediaFromReaderWithSize(
	ctx context.Context,
	uploadType schemes.UploadType,
	reader io.Reader,
	fileName string,
	fileSize int64,
	result any,
) error {
	endpoint, err := a.getUploadURL(ctx, uploadType)
	if err != nil {
		return err
	}

	fileName = multipartFileName(fileName)
	contentType, contentLength, boundary, err := multipartEnvelope(fileName, fileSize)
	if err != nil {
		return fmt.Errorf("prepare multipart: %w", err)
	}

	bodyReader, bodyWriter := io.Pipe()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.Url, bodyReader)
	if err != nil {
		_ = bodyReader.Close()
		_ = bodyWriter.Close()
		return fmt.Errorf("create upload request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	req.ContentLength = contentLength

	go func() {
		writer := multipart.NewWriter(bodyWriter)
		if err := writer.SetBoundary(boundary); err != nil {
			_ = bodyWriter.CloseWithError(fmt.Errorf("set multipart boundary: %w", err))
			return
		}

		fileWriter, err := writer.CreateFormFile("data", fileName)
		if err != nil {
			_ = bodyWriter.CloseWithError(fmt.Errorf("create form file: %w", err))
			return
		}
		if _, err = io.Copy(fileWriter, reader); err != nil {
			_ = bodyWriter.CloseWithError(fmt.Errorf("copy file data: %w", err))
			return
		}
		if err = writer.Close(); err != nil {
			_ = bodyWriter.CloseWithError(fmt.Errorf("close multipart writer: %w", err))
			return
		}
		_ = bodyWriter.Close()
	}()

	resp, err := a.client.do(req)
	if err != nil {
		return fmt.Errorf("upload: %w", err)
	}
	defer a.client.closer("uploadMediaFromReaderWithSize body", resp.Body)

	return a.decodeUploadResponse(resp, uploadType, endpoint.Token, result)
}

func (a *uploads) decodeUploadResponse(resp *http.Response, uploadType schemes.UploadType, token string, result any) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &schemes.Error{}
		if decodeErr := jsoniter.NewDecoder(resp.Body).Decode(apiErr); decodeErr == nil {
			return &APIError{
				Code:    resp.StatusCode,
				Message: apiErr.Code,
				Details: apiErr.Message,
			}
		}
		return fmt.Errorf("upload: HTTP %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if uploadType == schemes.AUDIO || uploadType == schemes.VIDEO {
		if info, ok := result.(*schemes.UploadedInfo); ok {
			info.Token = token
			return nil
		}
	}

	if err := jsoniter.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}

	return nil
}

func multipartEnvelope(fileName string, fileSize int64) (string, int64, string, error) {
	header := &bytes.Buffer{}
	writer := multipart.NewWriter(header)
	boundary := writer.Boundary()

	if _, err := writer.CreateFormFile("data", fileName); err != nil {
		return "", 0, "", err
	}

	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return "", 0, "", err
	}

	return contentType, int64(header.Len()) + fileSize, boundary, nil
}

func multipartFileName(fileName string) string {
	if fileName == "" {
		return "file"
	}

	return filepath.Base(fileName)
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
