package core

import (
	"context"
)

//go:generate mockgen -source=ports.go -destination=mocks/mock.go

type DB interface {
	GetComicsIdsByWords(ctx context.Context, words []string) ([]int, error)
	GetComicsesInfoByIds(ctx context.Context, ids []int) ([]ComicsInfo, error)
	GetWords2Ids(ctx context.Context) (map[string]*[]int, error)
}

type Index interface {
	Update(newVals map[string]*[]int)
	GetComicsIdsByWords(ctx context.Context, words []string) ([]int, error)
}

type Words interface {
	Norm(ctx context.Context, phrase string) ([]string, error)
}

type Searcher interface {
	Search(ctx context.Context, limit int, phrase string) ([]Comics, error)
	SearchIndex(ctx context.Context, limit int, phrase string) ([]Comics, error)
	SearchRandom(ctx context.Context, limit int) ([]Comics, error)
	SubscribeOnUpdateIndex(updateEventName string) error
	SubscribeInitiatorIndexUpdate()
	UpdateIndex()
}

type Initiator interface {
	UpdateIndex(searcher Searcher)
}

type UpdateSubscriber interface {
	SubscribeUpdateEvent(updateEventName string, searcher Searcher) error
}
