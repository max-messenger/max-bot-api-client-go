package model

type Subscription struct {
	URL            string   `json:"url"`
	SelfSignedCert string   `json:"self_signed_cert"`
	Time           int64    `json:"time"`
	UpdateTypes    []string `json:"update_types"`
	Version        string   `json:"version"`
}

type GetSubscriptionsResult struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type SubscriptionRequestBody struct {
	URL            string   `json:"url"`
	Secret         string   `json:"secret,omitempty"`
	SelfSignedCert string   `json:"self_signed_cert,omitempty"`
	UpdateTypes    []string `json:"update_types,omitempty"`
	Version        string   `json:"version,omitempty"`
}
