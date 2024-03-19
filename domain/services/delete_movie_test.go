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

func TestMovieService__DeleteMovie(t *testing.T) {
	movie := randomMovie()

	testCases := []struct {
		name          string
		req           *pb.MovieIDParams
		buildStubs    func(store *mock.MockStore)
		checkResponse func(t *testing.T, res *pb.MovieResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.MovieIDParams{
				Id: movie.ID,
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					DeleteMovie(gomock.Any(), gomock.Eq(movie.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.MovieResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "NotFound",
			req: &pb.MovieIDParams{
				Id: 0,
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					DeleteMovie(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.ErrRecordNotFound)
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
			req:  &pb.MovieIDParams{},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					DeleteMovie(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.ErrInternalError)
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
			res, err := movieServer.DeleteMovie(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
