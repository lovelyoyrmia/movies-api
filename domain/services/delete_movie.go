package services

import (
	"context"
	"errors"

	"github.com/lovelyoyrmia/movies-api/domain/pb"
	"github.com/lovelyoyrmia/movies-api/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *MovieService) DeleteMovie(ctx context.Context, req *pb.MovieIDParams) (*pb.MovieResponse, error) {
	err := service.store.DeleteMovie(ctx, req.GetId())

	if errors.Is(err, db.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "movie not found")
	}

	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.MovieResponse{
		Code:    codes.OK.String(),
		Message: "Successfully Deleted",
	}, nil
}
