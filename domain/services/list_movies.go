package services

import (
	"context"
	"time"

	"github.com/lovelyoyrmia/movies-api/domain/pb"
	"github.com/lovelyoyrmia/movies-api/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *MovieService) ListMovies(ctx context.Context, req *pb.ListMoviesParams) (*pb.ListMoviesResponse, error) {

	var limit int32

	if req.Limit == nil {
		limit = 100
	} else {
		limit = req.GetLimit()
	}

	movies, err := service.store.ListMovies(ctx, limit)

	if err != nil {
		return nil, status.Error(codes.Aborted, db.ErrInternalError.Error())
	}

	var moviesRes []*pb.Movie

	for _, v := range movies {
		moviesRes = append(moviesRes, &pb.Movie{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			Image:       v.Image.String,
			Rating:      float32(v.Rating.Float64),
			CreatedAt:   v.CreatedAt.Time.Format(time.DateTime),
			UpdateAt:    v.UpdatedAt.Time.Format(time.DateTime),
		})
	}

	return &pb.ListMoviesResponse{
		Code:    codes.OK.String(),
		Message: "Success",
		Data:    moviesRes,
	}, nil
}
