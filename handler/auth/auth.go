package auth

import (
	"net/http"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/handler/request"
	"github.com/fox-one/pkg/logger"
)

func HandleAuthentication(session core.Session) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)

			user, err := session.Login(r)
			if err != nil {
				next.ServeHTTP(w, r)
				log.WithError(err).Debugln("api: guest access")
				return
			}

			next.ServeHTTP(w, r.WithContext(
				request.WithUser(ctx, user),
			))
		}

		return http.HandlerFunc(fn)
	}
}
