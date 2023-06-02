package middlewares

import (
	"errors"
	"log"
	"net/http"

	"github.com/leomarzochi/facebooklike/cmd/auth"
	"github.com/leomarzochi/facebooklike/cmd/helpers"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s, %s, %s", r.Method, r.Host, r.RequestURI)
		next(w, r)
	}
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.ValidateToken(r)
		if err != nil {
			helpers.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
