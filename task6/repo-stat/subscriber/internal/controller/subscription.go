package controller

import (
	"context"
	subscriberpb "repo-stat/proto/subscriber"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (sh *Handler) PostSubscription(ctx context.Context, request *subscriberpb.PostSubscriptionRequest) (*subscriberpb.PostSubscriptionResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "subscription is required")
	}

	err := sh.subscription.Create(ctx, request.Subscription.Owner, request.Subscription.Repo)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &subscriberpb.PostSubscriptionResponse{}, nil
}

func (h *Handler) ListSubscriptions(ctx context.Context, request *subscriberpb.ListSubscriptionRequest) (*subscriberpb.ListSubscriptionResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "subscription is required")
	}

	subscriptions, err := h.subscription.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := make([]*subscriberpb.Subscription, 0, len(subscriptions))
	for _, sub := range subscriptions {
		response = append(response, &subscriberpb.Subscription{
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	return &subscriberpb.ListSubscriptionResponse{Subscriptions: response}, nil
}

func (h *Handler) DeleteSubscription(ctx context.Context, request *subscriberpb.DeleteSubscriptionRequest) (*subscriberpb.DeleteSubscriptionResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "subscription is required")
	}

	err := h.subscription.Delete(ctx, request.Subscription.Owner, request.Subscription.Repo)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &subscriberpb.DeleteSubscriptionResponse{}, nil
}
