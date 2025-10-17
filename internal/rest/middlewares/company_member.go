package middlewares

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/restkit/token"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) CompanyMember(
	UserCtxKey interface{},
	allowedCompanyRoles map[string]bool,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, ok := ctx.Value(UserCtxKey).(token.UserData)
			if !ok {
				ape.RenderErr(w,
					problems.Unauthorized("Missing AuthorizationHeader header"),
				)

				return
			}

			placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
			if err != nil {
				ape.RenderErr(w,
					problems.BadRequest(validation.Errors{
						"place_id": err,
					})...,
				)

				return
			}

			place, err := s.place.Get(ctx, placeID, "en")
			if err != nil {
				s.log.Error("failed to get place", "error", err, "place_id", placeID)
				switch {
				case errors.Is(err, errx.ErrorPlaceNotFound):
					ape.RenderErr(w, problems.NotFound("Place not found"))
				default:
					ape.RenderErr(w, problems.InternalError())
				}

				return
			}

			if place.CompanyID == nil {
				s.log.Error("place has no associated company", "place_id", placeID)
				ape.RenderErr(w, problems.Forbidden("Place is not associated with any company"))
				return
			}

			if user.CompanyID == nil {
				s.log.Error("user has no associated company", "user_id", user.ID)
				ape.RenderErr(w, problems.Forbidden("User is not associated with any company"))
				return
			}

			if place.CompanyID != user.CompanyID {
				s.log.Error("user company does not match place company", "user_company_id", user.CompanyID, "place_company_id", place.CompanyID, "place_id", placeID)
				ape.RenderErr(w, problems.Forbidden("User does not belong to the company associated with the place"))
				return
			}

			if !(allowedCompanyRoles[user.Role]) {
				s.log.Error("user role not allowed", "user_id", user.ID, "user_role", user.Role)
				ape.RenderErr(w, problems.Forbidden("User does not have the required role"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
