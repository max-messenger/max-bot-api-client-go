package model

type Subscription struct {
	URL         string   `json:"url"`
	Time        int64    `json:"time"`
	UpdateTypes []string `json:"update_types"`
}

type GetSubscriptionsResult struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type SubscriptionRequestBody struct {
	URL         string   `json:"url"`
	Secret      string   `json:"secret,omitempty"`
	UpdateTypes []string `json:"update_types,omitempty"`
}
