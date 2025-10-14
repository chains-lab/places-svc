package middlewares

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/restkit/auth"
)

func (s Service) CompanyRoleGrant(
	UserCtxKey interface{},
	allowedCompanyRoles map[string]bool,
	allowedSysadminRoles map[string]bool,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, ok := ctx.Value(UserCtxKey).(auth.UserData)
			if !ok {
				ape.RenderErr(w,
					problems.Unauthorized("Missing AuthorizationHeader header"),
				)

				return
			}

			_ = user

			//TODO: implement company employee role check

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
