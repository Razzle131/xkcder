package db

import (
	"context"
	"errors"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"yadro.com/course/search/core"
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

// ID не уникальны, повторы ID указывают на наличие нескольких слов из запроса
func (db *DB) GetComicsIdsByWords(ctx context.Context, words []string) ([]int, error) {
	query := `SELECT comics_id FROM key_words WHERE key_word = ANY($1);`

	var res []int
	if err := db.conn.SelectContext(ctx, &res, query, words); err != nil {
		db.log.Error("db get comics ids", "error", err)
		return nil, errors.New("failed to get comics id")
	}

	if len(res) == 0 {
		return nil, core.ErrNotFound
	}

	return res, nil
}

func (db *DB) GetComicsesInfoByIds(ctx context.Context, ids []int) ([]core.ComicsInfo, error) {
	query := `SELECT comics_id ID, image_url URL, words_total PhraseLen FROM comics WHERE comics_id = ANY($1)`

	var res []core.ComicsInfo
	if err := db.conn.SelectContext(ctx, &res, query, ids); err != nil {
		db.log.Error("db get comics info by id", "error", err)
		return nil, errors.New("failed to get comics by id")
	}

	if len(res) == 0 {
		return nil, core.ErrNotFound
	}

	return res, nil
}

func (db *DB) GetWords2Ids(ctx context.Context) (map[string]*[]int, error) {
	query := `SELECT key_word, comics_id FROM key_words;`

	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		db.log.Error("db get words to ids", "error", err)
		return nil, errors.New("failed to get key words and ids")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			db.log.Error("close rows", "error", err)
		}
	}()

	res := make(map[string]*[]int)
	for rows.Next() {
		var keyWord string
		var id int

		if err := rows.Scan(&keyWord, &id); err != nil {
			db.log.Error("db parse rows", "error", err)
			return nil, errors.New("failed to get key words and ids")
		}

		if res[keyWord] != nil {
			*res[keyWord] = append(*res[keyWord], id)
			continue
		}

		res[keyWord] = &[]int{id}
	}

	return res, nil
}
