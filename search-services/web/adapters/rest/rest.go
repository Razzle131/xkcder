package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"text/template"

	"yadro.com/course/web/core"
)

const tokenHeader = "Authorization"

func NewLogoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./html/logo2.png")
	}
}

func NewIndexHandler(logger *slog.Logger, indexPath string) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(indexPath))
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, nil)
		if err != nil {
			logger.Error("execute template", "error", err)
		}
	}
}

func NewSearchHandler(logger *slog.Logger, searchPath string, api core.Api) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(searchPath))
	return func(w http.ResponseWriter, r *http.Request) {
		phrase := r.URL.Query().Get("phrase")
		if len(phrase) == 0 {
			phrase = "empty"
		}

		resp, err := api.Search(phrase)
		if err != nil {
			logger.Error("api search", "error", err)
		}

		err = tmpl.Execute(w, resp)
		if err != nil {
			logger.Error("execute template", "error", err)
		}
	}
}

func NewAdminHandler(logger *slog.Logger, adminPath, tokenCookie string, api core.Api) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(adminPath))
	return func(w http.ResponseWriter, r *http.Request) {
		needsAuthorization := true
		defer func() {
			err := tmpl.Execute(w, needsAuthorization)
			if err != nil {
				logger.Error("execute template", "error", err)
			}
		}()

		token, err := r.Cookie(tokenCookie)
		if err != nil {
			return
		}

		err = api.Verify(tokenHeader, "Token "+token.Value)
		if err != nil {
			return
		}

		needsAuthorization = false
	}
}

func NewLoginHandler(tokenCookie string, api core.Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req core.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := api.Login(req.Name, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		cookie := http.Cookie{
			Name:  tokenCookie,
			Value: string(token),
		}

		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	}
}

func NewUpdateHandler(tokenCookie string, api core.Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(tokenCookie)
		if err != nil {
			return
		}

		err = api.DbUpdate(tokenHeader, "Token "+token.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewDropHandler(tokenCookie string, api core.Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(tokenCookie)
		if err != nil {
			return
		}

		err = api.DbDrop(tokenHeader, "Token "+token.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewStatsHandler(api core.Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := api.DbStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func NewSearchRandomHandler(searchPath string, api core.Api) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(searchPath))
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := api.SearchRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := tmpl.Execute(w, resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
