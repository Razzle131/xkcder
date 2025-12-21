package core

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

const maxRehandleTries = 5

type worker struct {
	xkcd  XKCD
	words Words
	db    DB
	log   *slog.Logger
}

func (w *worker) handle(ctx context.Context, comicsId int) error {
	info, err := w.xkcd.Get(ctx, comicsId)
	if err != nil {
		w.log.Error("xkcd get", "error", err, "comics_id", comicsId)
		return err
	}

	words, err := w.words.Norm(ctx, info.Title+" "+info.Description)
	if err != nil {
		w.log.Error("words norm", "error", err, "text", info.Title+" "+info.Description)
		return err
	}

	comics := Comics{
		ID:    info.ID,
		URL:   info.URL,
		Words: words,
	}

	err = w.db.Add(ctx, comics)
	if err != nil {
		w.log.Error("db add", "error", err, "comics", comics)
		return err
	}

	return nil
}

type workerPool struct {
	concurrency int
	pool        chan worker
	wg          sync.WaitGroup
	log         *slog.Logger
}

func newWorkerPool(concurrency int, log *slog.Logger) *workerPool {
	return &workerPool{
		concurrency: concurrency,
		pool:        make(chan worker, concurrency),
		wg:          sync.WaitGroup{},
		log:         log,
	}
}

func (wp *workerPool) create(xkcd XKCD, words Words, db DB, log *slog.Logger) {
	for range wp.concurrency {
		wp.pool <- worker{
			xkcd:  xkcd,
			words: words,
			db:    db,
			log:   log,
		}
	}
}

func (wp *workerPool) handle(ctx context.Context, comicsId, rehandleTry int) {
	w := <-wp.pool
	wp.wg.Go(func() {
		err := w.handle(ctx, comicsId)
		wp.pool <- w
		if err != nil && err != ErrNotFound && errors.Unwrap(err) != ErrDbAdapter && rehandleTry < maxRehandleTries {
			// +: повышение надежности сетевых запросов
			// -: при отвале word/xkcd некоторое время стучимся; также при ошибке бд не стучимся, т.к. либо бесполезно (отпала), либо нарушается ограничение и тп
			wp.log.Info("trying to rehandle comics", "comics_id", comicsId, "try", rehandleTry)
			wp.handle(ctx, comicsId, rehandleTry+1)
		}
	})
}

func (wp *workerPool) wait() {
	wp.wg.Wait()
	for range wp.concurrency {
		<-wp.pool
	}
}
