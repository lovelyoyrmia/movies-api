package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lovelyoyrmia/movies-api/internal/db"
	"github.com/lovelyoyrmia/movies-api/pkg/config"
	"github.com/lovelyoyrmia/movies-api/pkg/utils"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	database, err := db.NewDatabase(context.Background(), c)
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	store := db.NewStore(database.DB)

	seedMovies(store)
}

func seedMovies(store db.Store) {

	for i := 0; i < 30; i++ {
		randomMov := randomMovie()
		movie, err := store.AddMovie(context.Background(), db.AddMovieParams{
			Title:       randomMov.Title,
			Description: randomMov.Description,
			Rating:      randomMov.Rating,
			Image:       randomMov.Image,
		})

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("Movie Added : %v \n", movie.ID)
	}
}

func randomMovie() db.Movie {
	return db.Movie{
		ID:          int32(utils.RandomID()),
		Title:       utils.RandomTitle(),
		Description: utils.RandomDescription(),
		Image: pgtype.Text{
			String: utils.RandomString(15),
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
