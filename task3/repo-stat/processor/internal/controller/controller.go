package processor_controller

import (
	"context"
	processor_domain "repo-stat/processor/internal/domain"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProcessorService interface {
	GetRepoInfo(context.Context, string, string) (*processor_domain.RepoInfo, error)
	Ping(context.Context) (*processor_domain.Ping, error)
}

type processorController struct {
	ps ProcessorService
	processorpb.UnimplementedProcessorServer
}

func NewProcessorService(procserv ProcessorService) *processorController {
	return &processorController{
		ps: procserv,
	}
}

func (pc *processorController) GetInfoRepo(ctx context.Context, req *processorpb.GetInfoRepoRequest) (*processorpb.GetInfoRepoResponce, error) {
	resp, err := pc.ps.GetRepoInfo(ctx, req.GetOwner(), req.GetRepo())

	switch err {
	case processor_domain.BadRequest:
		return nil, status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	case processor_domain.InternalError:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	case processor_domain.ErrorNotFound:
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}

	return &processorpb.GetInfoRepoResponce{
		Fullname:    resp.FullName,
		Description: resp.Description,
		Forks:       resp.Forks,
		Stargazers:  resp.Stargazers,
		Createdat:   resp.CreatedAt,
	}, nil
}

func (pc *processorController) Ping(ctx context.Context, req *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	responce, _ := pc.ps.Ping(ctx)
	return &processorpb.PingResponse{
		Reply: responce.Reply,
	}, nil
}
