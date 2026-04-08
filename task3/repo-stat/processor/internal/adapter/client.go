package processor_adapter

import (
	"context"
	processor_domain "repo-stat/processor/internal/domain"
	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type collectorClient struct {
	cc   collectorpb.CollectorClient
	conn *grpc.ClientConn
}

func NewCollectorClient(address string) (*collectorClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, processor_domain.InternalError
	}

	return &collectorClient{
		cc:   collectorpb.NewCollectorClient(conn),
		conn: conn,
	}, nil

}

func (cc *collectorClient) GetRepoInfo(ctx context.Context, owner string, repo string) (*processor_domain.RepoInfo, error) {
	resp, err := cc.cc.GetInfoRepo(ctx, &collectorpb.GetInfoRepoRequest{
		Owner: owner,
		Repo:  repo,
	})

	if err != nil {
		status, _ := status.FromError(err)
		err = processor_domain.InternalError
		switch status.Code() {
		case codes.InvalidArgument:
			err = processor_domain.BadRequest
		case codes.NotFound:
			err = processor_domain.ErrorNotFound
		}
		return nil, err
	}

	return &processor_domain.RepoInfo{
		FullName:    resp.Fullname,
		Description: resp.Description,
		Forks:       resp.Forks,
		Stargazers:  resp.Stargazers,
		CreatedAt:   resp.Createdat,
	}, nil
}

func (cc *collectorClient) Ping(ctx context.Context) (*processor_domain.Ping, error) {
	return &processor_domain.Ping{Reply: "Pong"}, nil
}
