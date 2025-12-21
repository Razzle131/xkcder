package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"yadro.com/course/update/core"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func New(log *slog.Logger, address string) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	return &DB{
		log:  log,
		conn: db,
	}, nil
}

func (db *DB) CloseOrLog() {
	err := db.conn.Close()
	if err != nil {
		db.log.Error("close db conn", "error", err)
	}
}

func (db *DB) RollbackTxOrLog(tx *sqlx.Tx) {
	err := tx.Rollback()
	if err != nil {
		db.log.Error("rollback", "error", err)
	}
}

func (db *DB) Add(ctx context.Context, comics core.Comics) error {
	insertComics := `INSERT INTO comics (comics_id, image_url, words_total) VALUES ($1, $2, $3);`
	insertKeyWords := `INSERT INTO key_words (key_word, comics_id) VALUES ($1, $2);`

	tx, err := db.conn.BeginTxx(ctx, nil)
	if err != nil {
		db.log.Error("db.add begin tx", "error", err)
		return fmt.Errorf("%w: %v", core.ErrDbAdapter, err)
	}

	if _, err = tx.ExecContext(ctx, insertComics, comics.ID, comics.URL, len(comics.Words)); err != nil {
		db.log.Error("db.add insert comics", "error", err)
		db.RollbackTxOrLog(tx)
		return fmt.Errorf("%w: %v", core.ErrDbAdapter, err)
	}

	for _, word := range comics.Words {
		if _, err = tx.ExecContext(ctx, insertKeyWords, word, comics.ID); err != nil {
			db.log.Error("db.add insert words", "error", err)
			db.RollbackTxOrLog(tx)
			return fmt.Errorf("%w: %v", core.ErrDbAdapter, err)
		}
	}

	if err = tx.Commit(); err != nil {
		db.log.Error("db.add commit", "error", err)
		db.RollbackTxOrLog(tx)
		return fmt.Errorf("%w: %v", core.ErrDbAdapter, err)
	}

	return nil
}

func (db *DB) Stats(ctx context.Context) (core.DBStats, error) {
	query := `SELECT "total" WordsTotal, "unique" WordsUnique, "fetched" ComicsFetched FROM db_stats;`

	res := []core.DBStats{}
	if err := db.conn.SelectContext(ctx, &res, query); err != nil {
		db.log.Error("db.stats", "error", err)
		return core.DBStats{}, fmt.Errorf("%w: %v", core.ErrDbAdapter, err)
	}

	if len(res) > 1 {
		db.log.Error("db.stats bad result", "len", len(res))
		return core.DBStats{}, fmt.Errorf("%w: %v", core.ErrDbAdapter, errors.New("db stats bad result len"))
	}

	return res[0], nil
}

func (db *DB) IDs(ctx context.Context) ([]int, error) {
	query := `SELECT comics_id FROM comics;`

	res := make([]int, 0)
	if err := db.conn.SelectContext(ctx, &res, query); err != nil {
		db.log.Error("db stats", "error", err)
		return nil, fmt.Errorf("%w: %v", core.ErrDbAdapter, err)
	}

	slices.Sort(res)

	return res, nil
}

func (db *DB) Drop(ctx context.Context) error {
	dropKeyWords := `DELETE FROM key_words;`
	dropComics := `DELETE FROM comics;`

	tx, err := db.conn.BeginTxx(ctx, nil)
	if err != nil {
		db.log.Error("db.drop begin tx", "error", err)
		return fmt.Errorf("%w: %v", core.ErrDbAdapter, err)
	}

	if _, err = tx.ExecContext(ctx, dropKeyWords); err != nil {
		db.log.Error("db.drop drop words", "error", err)
		db.RollbackTxOrLog(tx)
		return fmt.Errorf("%w: %v", core.ErrDbAdapter, err)
	}

	if _, err = tx.ExecContext(ctx, dropComics); err != nil {
		db.log.Error("db.drop drop comics", "error", err)
		db.RollbackTxOrLog(tx)
		return fmt.Errorf("%w: %v", core.ErrDbAdapter, err)
	}
	if err = tx.Commit(); err != nil {
		db.log.Error("db.drop commit", "error", err)
		db.RollbackTxOrLog(tx)
		return fmt.Errorf("%w: %v", core.ErrDbAdapter, err)
	}

	return nil
}
