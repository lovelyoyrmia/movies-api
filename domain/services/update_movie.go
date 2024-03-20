package services

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lovelyoyrmia/movies-api/domain/pb"
	"github.com/lovelyoyrmia/movies-api/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *MovieService) UpdateMovie(ctx context.Context, req *pb.MovieParams) (*pb.MovieResponse, error) {
	movie, err := service.store.UpdateMovie(ctx, db.UpdateMovieParams{
		ID:          req.GetId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Rating: pgtype.Float8{
			Float64: float64(req.GetRating()),
			Valid:   true,
		},
		Image: pgtype.Text{
			String: req.GetImage(),
			Valid:  true,
		},
	})

	if errors.Is(err, db.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "movie not found")
	}

	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.MovieResponse{
		Code:    codes.OK.String(),
		Message: "Successfully Updated",
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
