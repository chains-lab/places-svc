package meta

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserData struct {
	ID        uuid.UUID `json:"sub,omitempty"`
	SessionID uuid.UUID `json:"session_id,omitempty"`
	Verified  bool      `json:"verified,omitempty"`
	Role      string    `json:"role,omitempty"`
}

func User(ctx context.Context) (UserData, error) {
	if ctx == nil {
		return UserData{}, status.Error(codes.Internal, "internal server error")
	}

	userData, ok := ctx.Value(UserCtxKey).(UserData)
	if !ok {
		return UserData{}, status.Error(codes.Unauthenticated, "missing metadata in request")
	}

	return userData, nil
}
