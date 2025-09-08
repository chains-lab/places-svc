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
	GetPlaceLocale(w http.ResponseWriter, r *http.Request)
	ListPlace(w http.ResponseWriter, r *http.Request)
	UpdatePlace(w http.ResponseWriter, r *http.Request)
	SetPlaceLocale(w http.ResponseWriter, r *http.Request)
	DeletePlace(w http.ResponseWriter, r *http.Request)
	SetTimetable(w http.ResponseWriter, r *http.Request)

	VerifyPlace(w http.ResponseWriter, r *http.Request)
	BlockPlace(w http.ResponseWriter, r *http.Request)

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
		r.Use(svc)
		r.Route("/v1", func(r chi.Router) {
			r.Route("/places", func(r chi.Router) {
				r.Get("/", h.ListPlace)
				r.With(auth).Post("/", h.CreatePlace)

				r.Route("/{place_id}", func(r chi.Router) {
					r.Get("/", h.GetPlace)
					r.With(auth).Put("/", h.UpdatePlace)
					r.With(auth).Delete("/", h.DeletePlace)

					r.Route("/{locale}", func(r chi.Router) {
						r.With(auth).Put("/", h.SetPlaceLocale)
						r.Get("/", h.GetPlaceLocale)
					})

					r.Route("/timetable", func(r chi.Router) {
						r.With(auth).Put("/", h.SetTimetable)
					})

					r.Post("/verified", h.VerifyPlace)
					r.Post("/block", h.BlockPlace)
				})
			})

			r.Route("/classes", func(r chi.Router) {
				r.Get("/", h.ListClass)
				r.With(auth).With(sysadmin).Post("/", h.CreateClass)

				r.Route("/{class_id}", func(r chi.Router) {
					r.Get("/", h.GetClass)
					r.With(auth).With(sysadmin).Put("/", h.UpdateClass)
					r.With(auth).With(sysadmin).Delete("/", h.DeleteClass)
				})
			})
		})
	})

	s.Start(ctx)

	<-ctx.Done()
	s.Stop(ctx)
}
