package domain

type Subscription struct {
	Owner string
	Repo  string
}

func NewSubscription(owner, repo string) *Subscription {
	return &Subscription{
		Owner: owner,
		Repo:  repo,
	}
}

type SubscriptionInfo struct {
	Repositories []Repository
}
