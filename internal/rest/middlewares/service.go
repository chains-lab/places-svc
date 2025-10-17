package middlewares

import (
	"context"
	"net/http"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/restkit/mdlv"
	"github.com/google/uuid"
)

type PlaceSvc interface {
	Get(ctx context.Context, placeID uuid.UUID, locale string) (models.Place, error)
}

type Service struct {
	place PlaceSvc
	log   logium.Logger
}

func New(log logium.Logger, svc PlaceSvc) Service {
	return Service{
		place: svc,
		log:   log,
	}
}

func (s Service) Auth(userCtxKey interface{}, skUser string) func(http.Handler) http.Handler {
	return mdlv.Auth(userCtxKey, skUser)
}

func (s Service) RoleGrant(userCtxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler {
	return mdlv.RoleGrant(userCtxKey, allowedRoles)
}
