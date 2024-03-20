package services

import (
	"context"
	"errors"
	"time"

	"github.com/lovelyoyrmia/movies-api/domain/pb"
	"github.com/lovelyoyrmia/movies-api/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *MovieService) DetailMovie(ctx context.Context, req *pb.MovieIDParams) (*pb.MovieResponse, error) {
	movie, err := service.store.DetailMovie(ctx, req.GetId())

	if errors.Is(err, db.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "movie not found")
	}

	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.MovieResponse{
		Code:    codes.OK.String(),
		Message: "Success",
		Data: &pb.Movie{
			Id:          movie.ID,
			Title:       movie.Title,
			Description: movie.Description,
			Image:       movie.Image.String,
			Rating:      float32(movie.Rating.Float64),
			CreatedAt:   movie.CreatedAt.Time.Format(time.DateTime),
			UpdatedAt:   movie.UpdatedAt.Time.Format(time.DateTime),
		},
	}, nil
}
