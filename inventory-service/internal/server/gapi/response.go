package gapi

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ymanshur/synasishouse/inventory/internal/typex"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	validationErrs         validation.Errors
	unprocessableEntityErr typex.UnProcessableEnity
	conflictErr            typex.Conflict
	notFoundErr            typex.NotFound
)

func translationError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.As(err, &validationErrs):
		return invalidArgumentError(convertValidationErrors(validationErrs))
	case errors.As(err, &unprocessableEntityErr):
		return status.Error(codes.InvalidArgument, unprocessableEntityErr.Error())
	case errors.As(err, &conflictErr):
		return status.Error(codes.AlreadyExists, conflictErr.Error())
	case errors.As(err, &notFoundErr):
		return status.Error(codes.NotFound, notFoundErr.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func convertValidationErrors(validationErrors validation.Errors) (violations []*errdetails.BadRequest_FieldViolation) {
	for field, err := range validationErrors {
		violations = append(violations, fieldViolation(field, err))
	}
	return
}

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")

	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}
