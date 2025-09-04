package problems

import (
	"context"
	"time"

	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/svc-errors/ape"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func nowRFC3339Nano() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

var ErrorInternal = ape.Declare("INTERNAL_ERROR")

func RaiseInternal(ctx context.Context, cause error) error {
	res, _ := status.New(codes.Internal, "internal server error").WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInternal.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInternal.Raise(
		cause,
		res,
	)
}

var ErrorPermissionDenied = ape.Declare("PERMISSION_DENIED")

func RaisePermissionDenied(ctx context.Context, cause error) error {
	res, _ := status.New(codes.PermissionDenied, cause.Error()).WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorPermissionDenied.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorPermissionDenied.Raise(
		cause,
		res,
	)
}

var ErrorUnauthenticated = ape.Declare("UNAUTHENTICATED")

func RaiseUnauthenticated(ctx context.Context, cause error) error {
	res, _ := status.New(codes.Unauthenticated, cause.Error()).WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorUnauthenticated.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorUnauthenticated.Raise(
		cause,
		res,
	)
}

var ErrorInvalidArgument = ape.Declare("INVALID_ARGUMENT")

func RaiseInvalidArgument(ctx context.Context, cause error, details ...*errdetails.BadRequest_FieldViolation) error {
	res, _ := status.New(codes.InvalidArgument, cause.Error()).WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInvalidArgument.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.BadRequest{
			FieldViolations: details,
		},
	)

	return ErrorInvalidArgument.Raise(
		cause,
		res,
	)
}
