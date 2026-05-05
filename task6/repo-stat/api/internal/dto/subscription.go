package dto

type SubscriptionResponse struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
}
