package services

import (
	"github.com/lovelyoyrmia/movies-api/domain/pb"
	"github.com/lovelyoyrmia/movies-api/internal/db"
)

type MovieService struct {
	pb.UnimplementedMovieServicesServer
	store db.Store
}

func NewMovieService(store db.Store) *MovieService {
	return &MovieService{
		store: store,
	}
}
