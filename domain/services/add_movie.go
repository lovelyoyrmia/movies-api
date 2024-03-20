package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lovelyoyrmia/movies-api/domain/pb"
	"github.com/lovelyoyrmia/movies-api/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *MovieService) AddMovie(ctx context.Context, req *pb.MovieParams) (*pb.MovieResponse, error) {
	violations := validateMovieParamsRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	movie, err := service.store.AddMovie(ctx, db.AddMovieParams{
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

	if db.ErrorCode(err) == db.UniqueViolation {
		return nil, status.Error(codes.AlreadyExists, "movie already exists")
	}

	if err != nil {
		return nil, status.Error(codes.Aborted, db.ErrInternalError.Error())
	}

	return &pb.MovieResponse{
		Code:    codes.OK.String(),
		Message: "Successfully Added",
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
