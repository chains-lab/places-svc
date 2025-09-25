package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/places-svc/internal/api/rest/meta"
	"github.com/go-chi/chi/v5"
)

type Controller interface {

	// Places level controller
	CreatePlace(w http.ResponseWriter, r *http.Request)
	GetPlace(w http.ResponseWriter, r *http.Request)
	ListPlace(w http.ResponseWriter, r *http.Request)
	UpdatePlace(w http.ResponseWriter, r *http.Request)
	DeletePlace(w http.ResponseWriter, r *http.Request)
	SetTimetable(w http.ResponseWriter, r *http.Request)
	VerifyPlace(w http.ResponseWriter, r *http.Request)
	UnverifyPlace(w http.ResponseWriter, r *http.Request)
	ActivatePlace(w http.ResponseWriter, r *http.Request)
	DeactivatePlace(w http.ResponseWriter, r *http.Request)
	ListLocalesForPlace(w http.ResponseWriter, r *http.Request)
	SetLocalesForPlace(w http.ResponseWriter, r *http.Request)
	GetTimetable(w http.ResponseWriter, r *http.Request)
	DeleteTimetable(w http.ResponseWriter, r *http.Request)

	// Classes level controller
	CreateClass(w http.ResponseWriter, r *http.Request)
	GetClass(w http.ResponseWriter, r *http.Request)
	ListClass(w http.ResponseWriter, r *http.Request)
	UpdateClass(w http.ResponseWriter, r *http.Request)
	ActivateClass(w http.ResponseWriter, r *http.Request)
	DeactivateClass(w http.ResponseWriter, r *http.Request)
	DeleteClass(w http.ResponseWriter, r *http.Request)
}

func (s *Service) Run(ctx context.Context, c Controller) {
	svc := mdlv.ServiceGrant(enum.CitiesSVC, s.cfg.JWT.Service.SecretKey)
	auth := mdlv.Auth(meta.UserCtxKey, s.cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
	})
	moder := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
		roles.Moder: true,
	})

	s.router.Route("/places-svc/", func(r chi.Router) {
		r.Use(svc)
		r.Route("/v1", func(r chi.Router) {
			r.Route("/classes", func(r chi.Router) {
				r.Get("/", c.ListClass)
				r.With(auth).With(sysadmin).Post("/", c.CreateClass)

				r.Route("/{class_code}", func(r chi.Router) {
					r.Get("/", c.GetClass)
					r.With(auth).With(sysadmin).Put("/", c.UpdateClass)
					r.With(auth).With(sysadmin).Delete("/", c.DeleteClass)

					r.With(auth).With(sysadmin).Put("/activate", c.ActivateClass)
					r.With(auth).With(sysadmin).Put("/deactivate", c.DeactivateClass)
				})
			})

			r.Route("/places", func(r chi.Router) {
				r.Get("/", c.ListPlace)
				r.With(auth).Post("/", c.CreatePlace)

				r.Route("/{place_id}", func(r chi.Router) {
					r.Get("/", c.GetPlace)
					r.With(auth).Put("/", c.UpdatePlace)
					r.With(auth).Delete("/", c.DeletePlace)

					r.With(auth).Put("/activate", c.ActivatePlace)
					r.With(auth).Put("/deactivate", c.DeactivatePlace)
					r.With(auth).With(moder).Put("/verified", c.VerifyPlace)
					r.With(auth).With(moder).Put("/unverified", c.UnverifyPlace)

					r.Route("/locales", func(r chi.Router) {
						r.Get("/", c.ListLocalesForPlace)
						r.With(auth).Put("/", c.SetLocalesForPlace)
					})

					r.Route("/timetable", func(r chi.Router) {
						r.Get("/", c.GetTimetable)
						r.With(auth).Put("/", c.SetTimetable)
						r.With(auth).Delete("/", c.DeleteTimetable)
					})
				})
			})
		})
	})

	s.start(ctx)

	<-ctx.Done()
	s.stop(ctx)
}
