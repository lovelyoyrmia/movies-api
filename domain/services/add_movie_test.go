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

func TestMovieService__AddMovie(t *testing.T) {
	movie := randomMovie()

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
				Title:       movie.Title,
				Description: movie.Description,
				Image:       movie.Image.String,
				Rating:      float32(movie.Rating.Float64),
			},
			buildStubs: func(store *mock.MockStore) {
				arg := db.AddMovieParams{
					Title:       movie.Title,
					Description: movie.Description,
					Rating:      movie.Rating,
					Image:       movie.Image,
				}

				store.EXPECT().
					AddMovie(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(movie, nil)
			},
			checkResponse: func(t *testing.T, res *pb.MovieResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "AlreadyExists",
			req: &pb.MovieParams{
				Id:          movie.ID,
				Title:       movie.Title,
				Description: movie.Description,
				Image:       movie.Image.String,
				Rating:      float32(movie.Rating.Float64),
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					AddMovie(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Movie{}, db.ErrUniqueViolation)
			},
			checkResponse: func(t *testing.T, res *pb.MovieResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.AlreadyExists, st.Code())
			},
		},
		{
			name: "InvalidArgument",
			req:  &pb.MovieParams{},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					AddMovie(gomock.Any(), gomock.Any()).
					AnyTimes()

			},
			checkResponse: func(t *testing.T, res *pb.MovieResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "Aborted",
			req: &pb.MovieParams{
				Title:       movie.Title,
				Description: movie.Description,
				Image:       movie.Image.String,
				Rating:      float32(movie.Rating.Float64),
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					AddMovie(gomock.Any(), gomock.Any()).
					AnyTimes().
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
			res, err := movieServer.AddMovie(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
