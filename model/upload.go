package model

type UploadEndpoint struct {
	Token string `json:"token,omitempty"`
	Url   string `json:"url"`
}
