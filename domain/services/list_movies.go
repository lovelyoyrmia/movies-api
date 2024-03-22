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
	
	var newMovies []db.Movie

	if req.Title == nil {
		movies, err := service.store.ListMovies(ctx, limit)
		if err != nil {
			return nil, status.Error(codes.Aborted, db.ErrInternalError.Error())
		}
		newMovies = movies
	} else {
		movies, err := service.store.ListMoviesByTitle(ctx, db.ListMoviesByTitleParams{
			Title: req.GetTitle(),
			Limit: limit,
		})
		if err != nil {
			return nil, status.Error(codes.Aborted, db.ErrInternalError.Error())
		}
		newMovies = movies
	}

	var moviesRes []*pb.Movie

	for _, v := range newMovies {
		moviesRes = append(moviesRes, &pb.Movie{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			Image:       v.Image.String,
			Rating:      float32(v.Rating.Float64),
			CreatedAt:   v.CreatedAt.Time.Format(time.DateTime),
			UpdatedAt:   v.UpdatedAt.Time.Format(time.DateTime),
		})
	}

	return &pb.ListMoviesResponse{
		Code:    codes.OK.String(),
		Message: "Success",
		Data:    moviesRes,
	}, nil
}
