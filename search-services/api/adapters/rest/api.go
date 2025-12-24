package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"yadro.com/course/api/core"
)

type PingResponse struct {
	Replies map[string]string `json:"replies"`
}

func NewPingHandler(log *slog.Logger, pingers map[string]core.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := PingResponse{
			Replies: core.PingServices(r.Context(), pingers),
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

func NewUpdateHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := updater.Update(r.Context())
		if err == core.ErrAlreadyExists {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("updater error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type UpdateStatsResponse struct {
	WordsTotal    int `json:"words_total"`
	WordsUnique   int `json:"words_unique"`
	ComicsFetched int `json:"comics_fetched"`
	ComicsTotal   int `json:"comics_total"`
}

func NewUpdateStatsHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := updater.Stats(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("updater error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		resp := UpdateStatsResponse{
			WordsTotal:    stats.WordsTotal,
			WordsUnique:   stats.WordsUnique,
			ComicsFetched: stats.ComicsFetched,
			ComicsTotal:   stats.ComicsTotal,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

type UpdateStatusResponse struct {
	Status string `json:"status"`
}

func NewUpdateStatusHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := updater.Status(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("updater error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		resp := UpdateStatusResponse{
			Status: string(status),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

func NewDropHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := updater.Drop(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("updater error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type Comics struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type SearchResponse struct {
	Comics []Comics `json:"comics"`
	Total  int      `json:"total"`
}

const (
	queryLimitParamName  = "limit"
	queryPhraseParamName = "phrase"
)

func parseQuerySearchParams(queryLimit, queryPhrase string) (int, string, error) {
	limit := 10

	if len(queryLimit) > 0 {
		var err error
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			return 0, "", errors.New("failed to parse limit parameter")
		}
	}

	if limit <= 0 {
		return 0, "", errors.New("bad limit parameter")
	}

	if len(queryPhrase) == 0 {
		return 0, "", errors.New("bad phrase parameter")
	}

	return limit, queryPhrase, nil
}

func NewSearchHandler(log *slog.Logger, searcher core.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, phrase, err := parseQuerySearchParams(r.URL.Query().Get(queryLimitParamName), r.URL.Query().Get(queryPhraseParamName))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := searcher.Search(r.Context(), phrase, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res2 := make([]Comics, 0, len(res))
		for _, c := range res {
			res2 = append(res2, Comics{
				Id:  c.ID,
				Url: c.URL,
			})
		}

		if err := json.NewEncoder(w).Encode(SearchResponse{Comics: res2, Total: len(res2)}); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

func NewIndexSearchHandler(log *slog.Logger, searcher core.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, phrase, err := parseQuerySearchParams(r.URL.Query().Get("limit"), r.URL.Query().Get("phrase"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := searcher.SearchIndex(r.Context(), phrase, limit)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res2 := make([]Comics, 0, len(res))
		for _, c := range res {
			res2 = append(res2, Comics{
				Id:  c.ID,
				Url: c.URL,
			})
		}

		if err := json.NewEncoder(w).Encode(SearchResponse{Comics: res2, Total: len(res2)}); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewLoginHandler(log *slog.Logger, aaa core.AAA) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creditionals LoginRequest
		err := json.NewDecoder(r.Body).Decode(&creditionals)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token, err := aaa.Login(creditionals.Name, creditionals.Password)
		if err == core.ErrNotAuthorized {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write([]byte(token)); err != nil {
			log.Error("cannot write reply", "error", err)
		}
	}
}

func NewRandomComicsHandler(log *slog.Logger, searcher core.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, _, err := parseQuerySearchParams(r.URL.Query().Get("limit"), "placeholder")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := searcher.SearchRandom(r.Context(), limit)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res2 := make([]Comics, 0, len(res))
		for _, c := range res {
			res2 = append(res2, Comics{
				Id:  c.ID,
				Url: c.URL,
			})
		}

		if err := json.NewEncoder(w).Encode(SearchResponse{Comics: res2, Total: len(res2)}); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}
