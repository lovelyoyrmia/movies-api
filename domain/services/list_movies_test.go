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

func TestMovieService__ListMovies(t *testing.T) {
	var movies []db.Movie

	limit := int32(10)

	for i := 1; i <= int(limit); i++ {
		movies = append(movies, randomMovie())
	}

	testCases := []struct {
		name          string
		req           *pb.ListMoviesParams
		buildStubs    func(store *mock.MockStore)
		checkResponse func(t *testing.T, res *pb.ListMoviesResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.ListMoviesParams{
				Limit: &limit,
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					ListMovies(gomock.Any(), gomock.Eq(limit)).
					Times(1).
					Return(movies, nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListMoviesResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "Aborted",
			req:  &pb.ListMoviesParams{},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					ListMovies(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Movie{}, db.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.ListMoviesResponse, err error) {
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
			res, err := movieServer.ListMovies(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
