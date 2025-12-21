package middleware

import (
	"net/http"
)

func Concurrency(next http.HandlerFunc, limit int) http.HandlerFunc {
	sem := make(chan struct{}, limit)
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case sem <- struct{}{}:
		default:
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		next(w, r)
		<-sem
	}
}
