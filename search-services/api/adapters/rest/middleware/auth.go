package middleware

import (
	"net/http"
	"regexp"
	"strings"
)

const (
	tokenHeader = "Authorization"
	tokenRegex  = "Token\\s.*"
)

//go:generate mockgen -source=auth.go -destination=mocks/mock.go

type TokenVerifier interface {
	Verify(token string) error
}

func Auth(next http.HandlerFunc, verifier TokenVerifier) http.HandlerFunc {
	regex := regexp.MustCompile(tokenRegex)
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(tokenHeader)
		if !regex.Match([]byte(token)) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		splitedToken := strings.Split(token, " ")[1]

		err := verifier.Verify(splitedToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
