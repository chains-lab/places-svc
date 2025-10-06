package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal"
	"github.com/chains-lab/places-svc/internal/rest/meta"
	"github.com/go-chi/chi/v5"
)

type Controller interface {

	// Places level controller
	CreatePlace(w http.ResponseWriter, r *http.Request)

	GetPlace(w http.ResponseWriter, r *http.Request)
	FilterPlace(w http.ResponseWriter, r *http.Request)

	UpdatePlace(w http.ResponseWriter, r *http.Request)
	UpdateVerifyPlace(w http.ResponseWriter, r *http.Request)
	UpdatePlaceStatus(w http.ResponseWriter, r *http.Request)

	DeletePlace(w http.ResponseWriter, r *http.Request)

	SetTimetable(w http.ResponseWriter, r *http.Request)
	GetTimetable(w http.ResponseWriter, r *http.Request)
	DeleteTimetable(w http.ResponseWriter, r *http.Request)

	SetLocalesForPlace(w http.ResponseWriter, r *http.Request)
	GetLocalesForPlace(w http.ResponseWriter, r *http.Request)

	CreateClass(w http.ResponseWriter, r *http.Request)

	GetClass(w http.ResponseWriter, r *http.Request)
	FilterClass(w http.ResponseWriter, r *http.Request)

	UpdateClass(w http.ResponseWriter, r *http.Request)
	ActivateClass(w http.ResponseWriter, r *http.Request)
	DeactivateClass(w http.ResponseWriter, r *http.Request)

	DeleteClass(w http.ResponseWriter, r *http.Request)
}

type Middleware interface {
	CompanyRoleGrant(
		UserCtxKey interface{},
		allowedCompanyRoles map[string]bool,
		allowedSysadminRoles map[string]bool,
	) func(http.Handler) http.Handler
}

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, m Middleware, c Controller) {
	svc := mdlv.ServiceGrant(cfg.Service.Name, cfg.JWT.Service.SecretKey)
	auth := mdlv.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
	})
	sysmoder := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
		roles.Moder: true,
	})

	companyAdmin := m.CompanyRoleGrant(meta.UserCtxKey, map[string]bool{
		"owner": true,
		"admin": true,
	}, map[string]bool{
		roles.Admin: true,
	})

	companyModer := m.CompanyRoleGrant(meta.UserCtxKey, map[string]bool{
		"owner": true,
		"admin": true,
		"moder": true,
	}, map[string]bool{
		roles.Admin: true,
	})

	r := chi.NewRouter()

	r.Route("/places-svc/", func(r chi.Router) {
		r.Use(svc)
		r.Route("/v1", func(r chi.Router) {
			r.Route("/classes", func(r chi.Router) {
				r.Get("/", c.FilterClass)
				r.With(auth, sysadmin).Post("/", c.CreateClass)

				r.Route("/{class_code}", func(r chi.Router) {
					r.Get("/", c.GetClass)

					r.Group(func(r chi.Router) {
						r.Use(auth, sysadmin)
						r.Put("/", c.UpdateClass)
						r.Delete("/", c.DeleteClass)

						r.Put("/activate", c.ActivateClass)
						r.Put("/deactivate", c.DeactivateClass)
					})
				})
			})

			r.Route("/places", func(r chi.Router) {
				r.Get("/", c.FilterPlace)
				r.With(auth).Post("/", c.CreatePlace)

				r.Route("/{place_id}", func(r chi.Router) {
					r.Get("/", c.GetPlace)

					r.Route("/locales", func(r chi.Router) {
						r.Get("/", c.GetLocalesForPlace)

						r.With(auth, companyModer).Put("/", c.SetLocalesForPlace)
					})

					r.Route("/timetable", func(r chi.Router) {
						r.Get("/", c.GetTimetable)

						r.Group(func(r chi.Router) {
							r.Use(auth, companyModer)
							r.Put("/", c.SetTimetable)
							r.Delete("/", c.DeleteTimetable)
						})
					})

					r.Group(func(r chi.Router) {
						r.Use(auth)

						r.With(companyModer).Put("/", c.UpdatePlace)
						r.With(companyAdmin).Delete("/", c.DeletePlace)

						r.With(companyAdmin).Put("/status", c.UpdatePlaceStatus)
					})

					r.Group(func(r chi.Router) {
						r.Use(auth)

						r.With(sysmoder).Put("/verify", c.UpdateVerifyPlace)
					})
				})
			})
		})
	})

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	<-ctx.Done()

	log.Info("shutting down REST service")
}
