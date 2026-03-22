package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "collector/api/proto"
)

type CollectorClient struct {
	conn   *grpc.ClientConn
	client pb.CollectorServiceClient
}

func NewCollectorClient(addr string) (*CollectorClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to collector: %w", err)
	}

	return &CollectorClient{
		conn:   conn,
		client: pb.NewCollectorServiceClient(conn),
	}, nil
}

func (c *CollectorClient) GetRepository(ctx context.Context, url string) (*pb.GetRepositoryResponse, error) {
	return c.client.GetRepositoryByURL(ctx, &pb.GetRepositoryByURLRequest{Url: url})
}

func (c *CollectorClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
