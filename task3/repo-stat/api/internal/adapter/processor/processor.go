package processor

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type processorClient struct {
	log *slog.Logger
	pc  processorpb.ProcessorClient
}

func NewProcessorClient(addres string, log *slog.Logger) (*processorClient, error) {

	conn, err := grpc.NewClient(addres, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := processorpb.NewProcessorClient(conn)

	return &processorClient{
		log: log,
		pc:  client,
	}, nil
}
func (pc *processorClient) Ping(ctx context.Context) domain.PingStatus {
	_, err := pc.pc.Ping(ctx, &processorpb.PingRequest{})
	if err != nil {
		pc.log.Error("processor ping failed", "error", err)
		return domain.PingStatusDown
	}
	return domain.PingStatusUp

}
func (pc *processorClient) GetInfoRepo(ctx context.Context, owner string, repo string) (*domain.RepoInfo, error) {
	resp, err := pc.pc.GetInfoRepo(ctx, &processorpb.GetInfoRepoRequest{
		Owner: owner,
		Repo:  repo,
	})

	if err != nil {
		pc.log.Error("processor getinforepo failed", "error", err)
		return nil, err
	}

	return &domain.RepoInfo{
		FullName:    resp.Fullname,
		Description: resp.Description,
		Forks:       resp.Forks,
		Stargazers:  resp.Stargazers,
		CreatedAt:   resp.Createdat,
	}, nil
}
