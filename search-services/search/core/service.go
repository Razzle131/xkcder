package core

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"math/rand"
	"slices"
)

// приколочено из-за лени(
const xkcdLastId = 3183

type Service struct {
	log       *slog.Logger
	db        DB
	words     Words
	index     Index
	initiator Initiator
	sub       UpdateSubscriber
}

func NewService(log *slog.Logger, db DB, words Words, index Index, initiator Initiator, sub UpdateSubscriber) *Service {
	s := &Service{
		log:       log,
		db:        db,
		words:     words,
		index:     index,
		initiator: initiator,
		sub:       sub,
	}

	return s
}

func (s *Service) SubscribeInitiatorIndexUpdate() {
	go s.initiator.UpdateIndex(s)
}

func (s *Service) SubscribeOnUpdateIndex(updateEventName string) error {
	return s.sub.SubscribeUpdateEvent(updateEventName, s)
}

func (s *Service) UpdateIndex() {
	newVals, err := s.db.GetWords2Ids(context.Background())
	if err != nil {
		s.log.Error("update index", "error", err)
		return
	}
	s.index.Update(newVals)
}

func (s *Service) rangeComicses(ctx context.Context, ids []int, limit int) ([]Comics, error) {
	countOccurrences := make(map[int]int)
	for _, id := range ids {
		if _, found := countOccurrences[id]; !found {
			countOccurrences[id] = 1
			continue
		}
		countOccurrences[id]++
	}

	info, err := s.db.GetComicsesInfoByIds(ctx, slices.Collect(maps.Keys(countOccurrences)))
	if err == ErrNotFound {
		return []Comics{}, nil
	}
	if err != nil {
		return nil, errors.New("failed to get comicses")
	}

	slices.SortFunc(info, func(a ComicsInfo, b ComicsInfo) int {
		if countOccurrences[a.ID] != countOccurrences[b.ID] {
			return countOccurrences[b.ID] - countOccurrences[a.ID]
		}

		aPercentage := float64(countOccurrences[a.ID]) / float64(a.PhraseLen)
		bPercentage := float64(countOccurrences[b.ID]) / float64(b.PhraseLen)
		if aPercentage > bPercentage {
			return -1
		}
		if aPercentage < bPercentage {
			return 1
		}
		return 0
	})

	res := make([]Comics, min(limit, len(info)))
	for i := range len(res) {
		res[i] = Comics{
			ID:  info[i].ID,
			URL: info[i].URL,
		}
	}

	return res, nil
}

func (s *Service) Search(ctx context.Context, limit int, phrase string) ([]Comics, error) {
	if limit <= 0 {
		return nil, errors.New("bad limit argument")
	}

	normalized, err := s.words.Norm(ctx, phrase)
	if err != nil {
		return nil, err
	}

	ids, err := s.db.GetComicsIdsByWords(ctx, normalized)
	if err == ErrNotFound {
		return []Comics{}, nil
	}
	if err != nil {
		return nil, errors.New("failed to get comicses")
	}

	return s.rangeComicses(ctx, ids, limit)
}

func (s *Service) SearchIndex(ctx context.Context, limit int, phrase string) ([]Comics, error) {
	if limit <= 0 {
		return nil, errors.New("bad limit argument")
	}

	normalized, err := s.words.Norm(ctx, phrase)
	if err != nil {
		return nil, err
	}

	ids, err := s.index.GetComicsIdsByWords(ctx, normalized)
	if err == ErrNotFound {
		return []Comics{}, nil
	}
	if err != nil {
		return nil, errors.New("failed to get comicses")
	}

	return s.rangeComicses(ctx, ids, limit)
}

func (s *Service) SearchRandom(ctx context.Context, limit int) ([]Comics, error) {
	if limit <= 0 {
		return nil, errors.New("bad limit argument")
	}

	ids := make([]int, 0, limit)
	for len(ids) < limit {
		chosen := rand.Intn(xkcdLastId) + 1
		if chosen != 404 {
			ids = append(ids, chosen)
		}
	}

	return s.rangeComicses(ctx, ids, limit)
}
