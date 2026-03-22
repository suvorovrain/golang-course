package grpc

import (
	"collector/internal/usecase/client"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "collector/api/proto"
	_ "collector/internal/usecase/client"
)

type Handler struct {
	pb.UnimplementedCollectorServiceServer
	useCase *client.UseCase
}

func NewHandler(useCase *client.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) GetRepositoryByURL(
	ctx context.Context,
	req *pb.GetRepositoryByURLRequest,
) (*pb.GetRepositoryResponse, error) {
	repo, err := h.useCase.Execute(ctx, req.Url)
	if err != nil {
		return h.mapError(err)
	}

	return &pb.GetRepositoryResponse{
		Repository: &pb.Repository{
			FullName:    repo.FullName,
			Description: repo.Description,
			Stars:       int32(repo.Stars),
			Forks:       int64(int32(repo.Forks)),
			CreatedAt:   repo.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Language:    repo.Language,
		},
		ErrorCode:    0,
		ErrorMessage: "",
	}, nil
}

func (h *Handler) mapError(err error) (*pb.GetRepositoryResponse, error) {
	var code codes.Code
	var message string

	switch {
	case err == client.ErrInvalidURL:
		code = codes.InvalidArgument
		message = "Invalid repository URL"
	case err == client.ErrRepoNotFound:
		code = codes.NotFound
		message = "Repository not found"
	default:
		code = codes.Internal
		message = "Internal server error"
	}

	return &pb.GetRepositoryResponse{
		ErrorCode:    int32(code),
		ErrorMessage: message,
	}, status.Error(code, message)
}
