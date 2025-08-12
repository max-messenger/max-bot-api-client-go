package maxbot

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
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
		return nil, err
	}
	defer fh.Close()
	return a.UploadMediaFromReader(ctx, uploadType, fh)
}

// UploadMediaFromUrl uploads the file from a remote server to the Max server.
func (a *uploads) UploadMediaFromUrl(ctx context.Context, uploadType schemes.UploadType, u url.URL) (*schemes.UploadedInfo, error) {
	respFile, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer respFile.Body.Close()
	return a.UploadMediaFromReader(ctx, uploadType, respFile.Body)
}

func (a *uploads) UploadMediaFromReader(ctx context.Context, uploadType schemes.UploadType, reader io.Reader) (*schemes.UploadedInfo, error) {
	result := new(schemes.UploadedInfo)
	return result, a.uploadMediaFromReader(ctx, uploadType, reader, result)
}

// UploadPhotoFromFile uploads photos to the Max server.
func (a *uploads) UploadPhotoFromFile(ctx context.Context, fileName string) (*schemes.PhotoTokens, error) {
	fh, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	result := new(schemes.PhotoTokens)
	return result, a.uploadMediaFromReader(ctx, schemes.PHOTO, fh, result)
}

// UploadPhotoFromBase64String uploads photos to the Max server.
func (a *uploads) UploadPhotoFromBase64String(ctx context.Context, code string) (*schemes.PhotoTokens, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(code))
	result := new(schemes.PhotoTokens)
	return result, a.uploadMediaFromReader(ctx, schemes.PHOTO, decoder, result)
}

// UploadPhotoFromUrl uploads the photo from a remote server to the Max server.
func (a *uploads) UploadPhotoFromUrl(ctx context.Context, url string) (*schemes.PhotoTokens, error) {
	respFile, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer respFile.Body.Close()
	result := new(schemes.PhotoTokens)
	return result, a.uploadMediaFromReader(ctx, schemes.PHOTO, respFile.Body, result)
}

// UploadPhotoFromReader uploads the photo from a reader.
func (a *uploads) UploadPhotoFromReader(ctx context.Context, reader io.Reader) (*schemes.PhotoTokens, error) {
	result := new(schemes.PhotoTokens)
	return result, a.uploadMediaFromReader(ctx, schemes.PHOTO, reader, result)
}

func (a *uploads) getUploadURL(ctx context.Context, uploadType schemes.UploadType) (*schemes.UploadEndpoint, error) {
	result := new(schemes.UploadEndpoint)
	values := url.Values{}
	values.Set(paramType, string(uploadType))
	body, err := a.client.request(ctx, http.MethodPost, "uploads", values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *uploads) uploadMediaFromReader(ctx context.Context, uploadType schemes.UploadType, reader io.Reader, result interface{}) error {
	endpoint, err := a.getUploadURL(ctx, uploadType)
	if err != nil {
		return err
	}
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("data", "file")
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, reader)
	if err != nil {
		return err
	}

	if err := bodyWriter.Close(); err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	if err := bodyWriter.Close(); err != nil {
		return err
	}

	resp, err := http.Post(endpoint.Url, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	switch result := result.(type) {
	case *schemes.UploadedInfo:
		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		result.Token = endpoint.Token
		// Print response body as string
		log.Println(string(body))
	default:
		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}
