package model

type ActionRequestBody struct {
	Action SenderAction `json:"action"`
}

type UploadedInfo struct {
	Token string `json:"token"`
}

type UserIdsList struct {
	UserIds []int64 `json:"user_ids"`
}
