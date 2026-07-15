package model

type InitData struct {
	Chat       ChatApp
	User       UserApp
	IP         string
	QueryID    string
	AuthDate   int64
	StartParam string
}

type ChatApp struct {
	Id   int64    `json:"id"`
	Type ChatType `json:"type"`
}

type UserApp struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	LanguageCode string `json:"language_code"`
	PhotoURL     string `json:"photo_url"`
}
