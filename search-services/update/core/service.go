package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync/atomic"
)

type Service struct {
	log         *slog.Logger
	db          DB
	xkcd        XKCD
	words       Words
	concurrency int
	wp          *workerPool
	isUpdating  atomic.Bool
}

func NewService(
	log *slog.Logger, db DB, xkcd XKCD, words Words, concurrency int,
) (*Service, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("wrong concurrency specified: %d", concurrency)
	}
	return &Service{
		log:         log,
		db:          db,
		xkcd:        xkcd,
		words:       words,
		concurrency: concurrency,
		wp:          newWorkerPool(concurrency, log),
		isUpdating:  atomic.Bool{},
	}, nil
}

func (s *Service) Update(ctx context.Context) (err error) {
	if !s.isUpdating.CompareAndSwap(false, true) {
		return ErrAlreadyUpdating
	}
	defer s.isUpdating.Store(false)

	s.wp.create(s.xkcd, s.words, s.db, s.log)

	xkcdLastId, err := s.xkcd.LastID(ctx)
	if err != nil {
		return err
	}

	ids, err := s.db.IDs(ctx)
	if err != nil {
		return errors.New("failed to get comics ids")
	}

	for id := 1; id <= xkcdLastId; id++ {
		if _, found := slices.BinarySearch(ids, id); found {
			continue
		}

		s.wp.handle(ctx, id, 0)
	}

	s.wp.wait()

	return nil
}

func (s *Service) Stats(ctx context.Context) (ServiceStats, error) {
	comicsTotal, err := s.xkcd.LastID(ctx)
	if err != nil {
		return ServiceStats{}, errors.New("failed to load xkcd stats")
	}

	dbStats, err := s.db.Stats(ctx)
	if err != nil {
		return ServiceStats{}, errors.New("failed to load db stats")
	}

	return ServiceStats{
		DBStats:     dbStats,
		ComicsTotal: comicsTotal - 1,
	}, nil
}

func (s *Service) Status(ctx context.Context) ServiceStatus {
	if s.isUpdating.Load() {
		return StatusRunning
	}
	return StatusIdle
}

func (s *Service) Drop(ctx context.Context) error {
	err := s.db.Drop(ctx)
	if err != nil {
		return errors.New("failed to clean tables")
	}

	return nil
}
