package model

type SimpleQueryResult struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}
