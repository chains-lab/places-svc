package interceptors

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/api/grpc/meta"
	"github.com/chains-lab/places-svc/internal/problems"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, problems.RaiseUnauthenticated(ctx, fmt.Errorf("no metadata found in incoming context"))
		}

		requestIDArr := md["x-request-id"]
		if len(requestIDArr) == 0 {
			return nil, problems.RaiseUnauthenticated(ctx, fmt.Errorf("request UserID not supplied"))
		}

		requestID, err := uuid.Parse(requestIDArr[0])
		if err != nil {
			return nil, problems.RaiseUnauthenticated(ctx, fmt.Errorf("invalid request UserID: %v", err))
		}

		ctx = context.WithValue(ctx, meta.RequestIDCtxKey, requestID)

		return handler(ctx, req)
	}
}
