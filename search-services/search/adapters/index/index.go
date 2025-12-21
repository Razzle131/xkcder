package index

import (
	"context"
	"sync"

	"yadro.com/course/search/core"
)

type Index struct {
	comics map[string]*[]int // keyword -> ids
	mu     sync.RWMutex
}

func New() *Index {
	return &Index{
		comics: make(map[string]*[]int),
		mu:     sync.RWMutex{},
	}
}

func (i *Index) Update(newVals map[string]*[]int) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.comics = newVals
}

func (i *Index) GetComicsIdsByWords(ctx context.Context, words []string) ([]int, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var res []int
	for _, word := range words {
		if i.comics[word] == nil {
			continue
		}

		res = append(res, *i.comics[word]...)
	}

	if len(res) == 0 {
		return nil, core.ErrNotFound
	}

	return res, nil
}
