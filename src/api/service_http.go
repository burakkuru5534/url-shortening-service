package api

import (
	"github.com/burakkuru5534/src/api/register"
	"github.com/burakkuru5534/src/api/sys"
	"github.com/burakkuru5534/src/auth"
	"github.com/burakkuru5534/src/helper"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
)

func HttpService() http.Handler {
	mux := chi.NewRouter()

	acors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	mux.Use(acors.Handler)

	mux.Route("/api", func(mr chi.Router) {
		mr.Group(func(r chi.Router) {
			//r.Get("/users/{id}", UserGet)
			r.Post("/login", sys.Login)
			r.Post("/register", register.NewRegister)

		})
		//protected end points
		mr.Group(func(r chi.Router) {
			//Token Middleware
			r.Use(jwtauth.Verifier(helper.Conf.Auth.JWTAuth))
			r.Use(jwtauth.Authenticator)
			r.Use(ProjectAuthenticator)

			r.Post("/url", UrlCreate)
			r.Get("/url", UrlList)

		})
	})

	return mux
}

func ProjectAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		tc := auth.TokenClaimsFromRequest(r)

		//TODO : if token.expire time.now ve db.isinvalid
		var isInvalid bool
		err = helper.App.DB.QueryRow("select is_invalid from logjwt where jwt = $1", jwtauth.TokenFromHeader(r)).Scan(&isInvalid)
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if (isInvalid) || (tc.ExpireIn < time.Now().Unix()) {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
