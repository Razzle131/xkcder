package core_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	mock_core "yadro.com/course/update/core/mocks"

	. "yadro.com/course/update/core"
)

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name          string
		dbBehaviour   func(db *mock_core.MockDB)
		wordBehaviour func(words *mock_core.MockWords)
		xkcdBehaviour func(xkcd *mock_core.MockXKCD)
		wantErr       bool
	}{
		{
			name: "OK",
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().IDs(t.Context()).Return([]int{}, nil)
				db.EXPECT().Add(t.Context(), Comics{
					ID:    1,
					Words: []string{"abc"},
				}).Return(nil)
			},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "abc ").Return([]string{"abc"}, nil)
			},
			xkcdBehaviour: func(xkcd *mock_core.MockXKCD) {
				xkcd.EXPECT().LastID(t.Context()).Return(1, nil)
				xkcd.EXPECT().Get(t.Context(), 1).Return(XKCDInfo{ID: 1, Title: "abc"}, nil)
			},
		},
		{
			name: "all comics in db already",
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().IDs(t.Context()).Return([]int{1}, nil)
			},
			wordBehaviour: func(words *mock_core.MockWords) {},
			xkcdBehaviour: func(xkcd *mock_core.MockXKCD) {
				xkcd.EXPECT().LastID(t.Context()).Return(1, nil)
			},
		},
		{
			name:          "XKCD error",
			dbBehaviour:   func(db *mock_core.MockDB) {},
			wordBehaviour: func(words *mock_core.MockWords) {},
			xkcdBehaviour: func(xkcd *mock_core.MockXKCD) {
				xkcd.EXPECT().LastID(t.Context()).Return(0, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "db error",
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().IDs(t.Context()).Return(nil, errors.New("some error"))
			},
			wordBehaviour: func(words *mock_core.MockWords) {},
			xkcdBehaviour: func(xkcd *mock_core.MockXKCD) {
				xkcd.EXPECT().LastID(t.Context()).Return(1, nil)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			db := mock_core.NewMockDB(c)
			words := mock_core.NewMockWords(c)
			xkcd := mock_core.NewMockXKCD(c)

			testCase.dbBehaviour(db)
			testCase.wordBehaviour(words)
			testCase.xkcdBehaviour(xkcd)

			service, err := NewService(slog.Default(), db, xkcd, words, 10)
			assert.NoError(t, err)

			err = service.Update(t.Context())

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestStats(t *testing.T) {
	testCases := []struct {
		name          string
		expected      ServiceStats
		dbBehaviour   func(db *mock_core.MockDB)
		xkcdBehaviour func(xkcd *mock_core.MockXKCD)
		wantErr       bool
	}{
		{
			name: "OK",
			expected: ServiceStats{
				DBStats: DBStats{
					WordsTotal:    1,
					WordsUnique:   1,
					ComicsFetched: 1,
				},
				ComicsTotal: 404,
			},
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().Stats(t.Context()).Return(DBStats{
					WordsTotal:    1,
					WordsUnique:   1,
					ComicsFetched: 1,
				}, nil)
			},
			xkcdBehaviour: func(xkcd *mock_core.MockXKCD) {
				xkcd.EXPECT().LastID(t.Context()).Return(405, nil)
			},
		},
		{
			name:        "xkcd error",
			dbBehaviour: func(db *mock_core.MockDB) {},
			xkcdBehaviour: func(xkcd *mock_core.MockXKCD) {
				xkcd.EXPECT().LastID(t.Context()).Return(0, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "db error",
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().Stats(t.Context()).Return(DBStats{}, errors.New("some error"))
			},
			xkcdBehaviour: func(xkcd *mock_core.MockXKCD) {
				xkcd.EXPECT().LastID(t.Context()).Return(405, nil)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			db := mock_core.NewMockDB(c)
			xkcd := mock_core.NewMockXKCD(c)

			testCase.dbBehaviour(db)
			testCase.xkcdBehaviour(xkcd)

			service, err := NewService(slog.Default(), db, xkcd, nil, 10)
			assert.NoError(t, err)

			res, err := service.Stats(t.Context())

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, res)
		})
	}
}
