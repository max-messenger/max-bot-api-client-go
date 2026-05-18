package maxbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type Upload struct {
	client *client
}

func newUpload(cli *client) *Upload {
	return &Upload{
		client: cli,
	}
}

func (u *Upload) Upload(ctx context.Context, uploadType model.UploadType, r io.Reader, name string, size int64) (token string, err error) {
	endpoint, err := u.getUploadURL(ctx, uploadType)
	if err != nil {
		return
	}

	name = multipartFileName(name)
	contentType, contentLength, boundary, err := multipartEnvelope(name, size)
	if err != nil {
		return
	}

	bodyReader, bodyWriter := io.Pipe()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.Url, bodyReader)
	if err != nil {
		_ = bodyReader.Close()
		_ = bodyWriter.Close()

		return
	}
	req.Header.Set("Content-Type", contentType)
	req.ContentLength = contentLength

	go func() {
		writer := multipart.NewWriter(bodyWriter)
		if gErr := writer.SetBoundary(boundary); gErr != nil {
			_ = bodyWriter.CloseWithError(fmt.Errorf("set multipart boundary: %w", gErr))

			return
		}

		fileWriter, gErr := writer.CreateFormFile(fieldData, name)
		if gErr != nil {
			_ = bodyWriter.CloseWithError(fmt.Errorf("create form file: %w", gErr))

			return
		}
		if _, gErr = io.Copy(fileWriter, r); gErr != nil {
			_ = bodyWriter.CloseWithError(fmt.Errorf("copy file data: %w", gErr))

			return
		}
		if gErr = writer.Close(); gErr != nil {
			_ = bodyWriter.CloseWithError(fmt.Errorf("close multipart writer: %w", gErr))

			return
		}
		_ = bodyWriter.Close()
	}()

	resp, err := u.client.do(req)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		err = parseResponseError(resp)

		return
	}

	if uploadType == model.UploadAudio || uploadType == model.UploadVideo {
		token = endpoint.Token

		return
	}

	result := model.UploadedInfo{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		err = fmt.Errorf("unmarshal response body: %w", err)

		return
	}

	token = result.Token

	return
}

func (u *Upload) getUploadURL(ctx context.Context, uploadType model.UploadType) (res model.UploadEndpoint, err error) {
	values := url.Values{}
	values.Set(paramType, string(uploadType))

	err = u.client.raw(ctx, http.MethodPost, pathUpload, values, nil, &res)
	if err != nil {
		err = fmt.Errorf("getUploadURL err: %w", err)

		return
	}

	return
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
		return defaultFileName
	}

	return filepath.Base(fileName)
}
