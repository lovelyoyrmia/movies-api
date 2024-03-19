package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lovelyoyrmia/movies-api/domain/pb"
	"github.com/lovelyoyrmia/movies-api/internal/db"
	"github.com/lovelyoyrmia/movies-api/internal/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMovieService__UpdateMovie(t *testing.T) {
	movie := randomMovie()
	newMovie := randomMovie()

	testCases := []struct {
		name          string
		req           *pb.MovieParams
		buildStubs    func(store *mock.MockStore)
		checkResponse func(t *testing.T, res *pb.MovieResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.MovieParams{
				Id:          movie.ID,
				Title:       newMovie.Title,
				Description: newMovie.Description,
				Image:       newMovie.Image.String,
				Rating:      float32(newMovie.Rating.Float64),
			},
			buildStubs: func(store *mock.MockStore) {
				arg := db.UpdateMovieParams{
					ID:          movie.ID,
					Title:       newMovie.Title,
					Description: newMovie.Description,
					Rating:      newMovie.Rating,
					Image:       newMovie.Image,
				}

				store.EXPECT().
					UpdateMovie(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(newMovie, nil)
			},
			checkResponse: func(t *testing.T, res *pb.MovieResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "NotFound",
			req: &pb.MovieParams{
				Id:          0,
				Title:       newMovie.Title,
				Description: newMovie.Description,
				Image:       newMovie.Image.String,
				Rating:      float32(newMovie.Rating.Float64),
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					UpdateMovie(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Movie{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.MovieResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "Aborted",
			req:  &pb.MovieParams{},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					UpdateMovie(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Movie{}, db.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.MovieResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Aborted, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			store := mock.NewMockStore(storeCtrl)

			tc.buildStubs(store)
			movieServer := NewMovieService(store)
			res, err := movieServer.UpdateMovie(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
