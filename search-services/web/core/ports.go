package core

type Api interface {
	DbDrop(tokenHeader string, token string) error
	DbStats() (UpdateStatsResponse, error)
	DbUpdate(tokenHeader string, token string) error
	Login(name string, password string) (string, error)
	Search(phrase string) (SearchResponse, error)
	SearchRandom() (SearchResponse, error)
	Verify(tokenHeader string, token string) error
}
