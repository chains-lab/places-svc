package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal"
	"github.com/chains-lab/places-svc/internal/rest/meta"
	"github.com/chains-lab/restkit/roles"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	CreatePlace(w http.ResponseWriter, r *http.Request)

	GetPlace(w http.ResponseWriter, r *http.Request)
	FilterPlace(w http.ResponseWriter, r *http.Request)

	UpdatePlace(w http.ResponseWriter, r *http.Request)
	UpdateVerifiedPlace(w http.ResponseWriter, r *http.Request)
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
	Auth(userCtxKey interface{}, skUser string) func(http.Handler) http.Handler
	RoleGrant(userCtxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler

	CompanyMember(
		UserCtxKey interface{},
		allowedCompanyRoles map[string]bool,
	) func(http.Handler) http.Handler

	CompanyMemberOrAdmin(
		UserCtxKey interface{},
		allowedCompanyRoles map[string]bool,
		allowedAdminRoles map[string]bool,
	) func(http.Handler) http.Handler
}

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, m Middleware, h Handlers) {
	auth := m.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)

	sysadmin := m.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
	})
	sysmoder := m.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
		roles.Moder: true,
	})

	companyAdmin := m.CompanyMember(meta.UserCtxKey, map[string]bool{
		"owner": true,
		"admin": true,
	})
	companyModer := m.CompanyMember(meta.UserCtxKey, map[string]bool{
		"owner": true,
		"admin": true,
		"moder": true,
	})

	r := chi.NewRouter()

	r.Route("/places-svc/", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/classes", func(r chi.Router) {
				r.Get("/", h.FilterClass)
				r.With(auth, sysadmin).Post("/", h.CreateClass)

				r.Route("/{class_code}", func(r chi.Router) {
					r.Get("/", h.GetClass)

					r.Group(func(r chi.Router) {
						r.Use(auth, sysadmin)
						r.Put("/", h.UpdateClass)
						r.Delete("/", h.DeleteClass)

						r.Put("/activate", h.ActivateClass)
						r.Put("/deactivate", h.DeactivateClass)
					})
				})
			})

			r.Route("/places", func(r chi.Router) {
				r.Get("/", h.FilterPlace)

				r.With(auth).Post("/", h.CreatePlace)
				r.With(auth, companyModer).Put("/", h.UpdatePlace)
				r.With(auth, companyAdmin).Delete("/", h.DeletePlace)

				r.Route("/{place_id}", func(r chi.Router) {
					r.Get("/", h.GetPlace)

					r.With(auth, sysmoder).Put("/verify", h.UpdateVerifiedPlace)
					r.With(auth, companyAdmin).Put("/status", h.UpdatePlaceStatus)

					r.Route("/locales", func(r chi.Router) {
						r.Get("/", h.GetLocalesForPlace)

						r.With(auth, companyModer).Put("/", h.SetLocalesForPlace)
					})

					r.Route("/timetable", func(r chi.Router) {
						r.Get("/", h.GetTimetable)

						r.Group(func(r chi.Router) {
							r.Use(auth, companyModer)
							r.Put("/", h.SetTimetable)
							r.Delete("/", h.DeleteTimetable)
						})
					})
				})
			})
		})
	})

	srv := &http.Server{
		Addr:              cfg.Rest.Port,
		Handler:           r,
		ReadTimeout:       cfg.Rest.Timeouts.Read,
		ReadHeaderTimeout: cfg.Rest.Timeouts.ReadHeader,
		WriteTimeout:      cfg.Rest.Timeouts.Write,
		IdleTimeout:       cfg.Rest.Timeouts.Idle,
	}

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down REST service...")
	case err := <-errCh:
		if err != nil {
			log.Errorf("REST server error: %v", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.Errorf("REST shutdown error: %v", err)
	} else {
		log.Info("REST server stopped")
	}
}
