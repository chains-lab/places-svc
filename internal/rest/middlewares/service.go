package middlewares

import (
	"net/http"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/restkit/mdlv"
)

type Service struct {
	log logium.Logger
}

func New(log logium.Logger) Service {
	return Service{
		log: log,
	}
}

func (s Service) Auth(userCtxKey interface{}, skUser string) func(http.Handler) http.Handler {
	return mdlv.Auth(userCtxKey, skUser)
}

func (s Service) RoleGrant(userCtxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler {
	return mdlv.RoleGrant(userCtxKey, allowedRoles)
}
