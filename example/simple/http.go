package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type httpClient struct {
	client *http.Client
}

func newHttpClient() *httpClient {
	return &httpClient{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	log.Printf("request: %s %s", req.Method, req.URL)
	response, err := c.client.Do(req)
	if err != nil {
		log.Printf("request error: %s", err)

		return nil, err
	}

	debugResponse(response)

	log.Printf("response: %s %d", response.Status, response.StatusCode)

	return response, nil
}

func debugResponse(response *http.Response) {
	// for debug
	bodyBytes, _ := io.ReadAll(response.Body)
	response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, bodyBytes, "", "\t")

	fmt.Println(prettyJSON.String())
}
