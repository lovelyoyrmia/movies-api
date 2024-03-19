package services

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lovelyoyrmia/movies-api/internal/db"
	"github.com/lovelyoyrmia/movies-api/pkg/utils"
)

func randomMovie() db.Movie {
	return db.Movie{
		ID:          int32(utils.RandomID()),
		Title:       utils.RandomString(),
		Description: utils.RandomString(),
		Image: pgtype.Text{
			String: utils.RandomString(),
			Valid:  true,
		},
		Rating: pgtype.Float8{
			Float64: utils.RandomRating(),
			Valid:   true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  time.Time{},
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Time{},
			Valid: true,
		},
	}
}
