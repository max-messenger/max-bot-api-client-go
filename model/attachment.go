package model

type Attachment struct {
	Type    AttachmentType `json:"type"`
	Payload Payload        `json:"payload"`

	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type Payload struct {
	Token     string      `json:"token,omitempty"`
	URL       string      `json:"url,omitempty"`
	Code      string      `json:"code,omitempty"`
	ContactID int64       `json:"contact_id,omitempty"` // for send contact
	Buttons   [][]*Button `json:"buttons,omitempty"`
}

type VideoUrls struct {
	Mp41080 *string `json:"mp4_1080"`
	Mp4720  *string `json:"mp4_720"`
	Mp4480  *string `json:"mp4_480"`
	Mp4360  *string `json:"mp4_360"`
	Mp4240  *string `json:"mp4_240"`
	Mp4144  *string `json:"mp4_144"`
	HLS     *string `json:"hls"`
}

type VideoAttachmentDetails struct {
	Token     string     `json:"token"`
	Urls      *VideoUrls `json:"urls,omitempty"`
	Thumbnail *Payload   `json:"thumbnail,omitempty"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Duration  int        `json:"duration"`
}
