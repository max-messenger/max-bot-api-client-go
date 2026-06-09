package model

type UploadEndpoint struct {
	Token string `json:"token,omitempty"`
	Url   string `json:"url"`
}

type PhotoTokens struct {
	Photos map[string]PhotoToken `json:"photos"`
}

type PhotoToken struct {
	Token string `json:"token"`
}
