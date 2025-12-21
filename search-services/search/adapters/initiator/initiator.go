package initiator

import (
	"time"

	"yadro.com/course/search/core"
)

type Initiator struct {
	ticker time.Ticker
}

func New(indexTTL time.Duration) *Initiator {
	return &Initiator{
		ticker: *time.NewTicker(indexTTL),
	}
}

func (i *Initiator) UpdateIndex(searcher core.Searcher) {
	for range i.ticker.C {
		searcher.UpdateIndex()
	}
}
