package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/places-svc/internal/api/rest/meta"
	"github.com/chains-lab/places-svc/internal/config"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {

	// Places level handlers
	CreatePlace(w http.ResponseWriter, r *http.Request)
	GetPlace(w http.ResponseWriter, r *http.Request)
	ListPlace(w http.ResponseWriter, r *http.Request)
	UpdatePlace(w http.ResponseWriter, r *http.Request)
	UpdatePlaceLocale(w http.ResponseWriter, r *http.Request)
	DeletePlace(w http.ResponseWriter, r *http.Request)
	SetTimetable(w http.ResponseWriter, r *http.Request)

	// Class level handlers
	CreateClass(w http.ResponseWriter, r *http.Request)
	GetClass(w http.ResponseWriter, r *http.Request)
	ListClass(w http.ResponseWriter, r *http.Request)
	UpdateClass(w http.ResponseWriter, r *http.Request)
	DeleteClass(w http.ResponseWriter, r *http.Request)
}

func (s *Service) Api(ctx context.Context, cfg config.Config, h Handlers) {
	svc := mdlv.ServiceGrant(constant.ServiceName, cfg.JWT.Service.SecretKey)
	auth := mdlv.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin:     true,
		roles.SuperUser: true,
	})

	s.router.Route("/places-svc/", func(r chi.Router) {

	})

	s.Start(ctx)

	<-ctx.Done()
	s.Stop(ctx)
}
