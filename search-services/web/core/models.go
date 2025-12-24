package core

type Comics struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type SearchResponse struct {
	Comics []Comics `json:"comics"`
	Total  int      `json:"total"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UpdateStatsResponse struct {
	WordsTotal    int `json:"words_total"`
	WordsUnique   int `json:"words_unique"`
	ComicsFetched int `json:"comics_fetched"`
	ComicsTotal   int `json:"comics_total"`
}
