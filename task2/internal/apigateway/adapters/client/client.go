package apigatewayclient

import (
	"context"
	"log"

	collectorpb "github.com/Chaice1/golang-course/task2/gen"
	apigatewaydomain "github.com/Chaice1/golang-course/task2/internal/apigateway/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type collectorClient struct {
	collectorpb.CollectorClient
}

func NewCollectorClient(cc collectorpb.CollectorClient) *collectorClient {
	return &collectorClient{CollectorClient: cc}
}

func (cc *collectorClient) GetRepoInfo(ctx context.Context, owner string, repo string) (*apigatewaydomain.RepoInfo, error) {

	repoinfo, err := cc.GetInfoRepo(ctx, &collectorpb.GetInfoRepoRequest{
		Owner: owner,
		Repo:  repo,
	})

	log.Println(repoinfo, err)
	if err == nil {
		return &apigatewaydomain.RepoInfo{
			FullName:    repoinfo.GetFullname(),
			Description: repoinfo.GetDescription(),
			Stargazers:  repoinfo.GetStargazers(),
			Forks:       repoinfo.GetForks(),
			CreatedAt:   repoinfo.GetCreatedat(),
		}, nil
	}

	status, _ := status.FromError(err)

	err = apigatewaydomain.InternalError
	switch status.Code() {
	case codes.InvalidArgument:
		err = apigatewaydomain.BadRequest
	case codes.NotFound:
		err = apigatewaydomain.NotFound
	}

	return nil, err

}
