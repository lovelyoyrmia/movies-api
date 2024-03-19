package services

import (
	"errors"
	"fmt"

	"github.com/lovelyoyrmia/movies-api/domain/pb"
	"github.com/lovelyoyrmia/movies-api/internal/db"
	"github.com/lovelyoyrmia/movies-api/pkg/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrEmptyFields = errors.New("cannot be empty")

type MovieService struct {
	pb.UnimplementedMovieServicesServer
	store db.Store
}

func NewMovieService(store db.Store) *MovieService {
	return &MovieService{
		store: store,
	}
}

func validateMovieParamsRequest(req *pb.MovieParams) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.Title == "" {
		err := fmt.Errorf("title %s", ErrEmptyFields.Error())
		return append(violations, validator.FieldViolation("title", err))
	}

	if req.Description == "" {
		err := fmt.Errorf("description %s", ErrEmptyFields.Error())
		return append(violations, validator.FieldViolation("description", err))
	}

	if err := validator.ValidateTitle(req.GetTitle()); err != nil {
		return append(violations, validator.FieldViolation("title", err))
	}

	return violations
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badReq := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid request")

	statusDetails, err := statusInvalid.WithDetails(badReq)
	if err != nil {
		return statusInvalid.Err()
	}
	return statusDetails.Err()
}
