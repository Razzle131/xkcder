package core_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	. "yadro.com/course/search/core"
	mock_core "yadro.com/course/search/core/mocks"
)

func TestSearch(t *testing.T) {
	testCases := []struct {
		name          string
		inputPhrase   string
		inputLimit    int
		expected      []Comics
		dbBehaviour   func(db *mock_core.MockDB)
		wordBehaviour func(words *mock_core.MockWords)
		wantErr       bool
	}{
		{
			name:        "OK",
			inputPhrase: "abcd",
			inputLimit:  1,
			expected:    []Comics{{1, ""}},
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().GetComicsIdsByWords(t.Context(), []string{"abcd"}).Return([]int{1}, nil)
				db.EXPECT().GetComicsesInfoByIds(t.Context(), []int{1}).Return([]ComicsInfo{{1, "", 4}}, nil)
			},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "abcd").Return([]string{"abcd"}, nil)
			},
		},
		{
			name:        "test ranging",
			inputPhrase: "abcd",
			inputLimit:  1,
			expected:    []Comics{{1, ""}},
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().GetComicsIdsByWords(t.Context(), []string{"abcd"}).Return([]int{1, 2}, nil)
				db.EXPECT().GetComicsesInfoByIds(t.Context(), []int{1, 2}).Return([]ComicsInfo{{1, "", 4}, {2, "", 8}}, nil)
			},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "abcd").Return([]string{"abcd"}, nil)
			},
		},
		{
			name:        "db error",
			inputPhrase: "abcd",
			inputLimit:  1,
			expected:    nil,
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().GetComicsIdsByWords(t.Context(), []string{"abcd"}).Return(nil, errors.New("some error"))
			},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "abcd").Return([]string{"abcd"}, nil)
			},
			wantErr: true,
		},
		{
			name:        "missing phrase",
			inputLimit:  1,
			dbBehaviour: func(db *mock_core.MockDB) {},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "").Return(nil, errors.New("some errror"))
			},
			wantErr: true,
		},
		{
			name:          "missing limit",
			inputPhrase:   "aaaa",
			dbBehaviour:   func(db *mock_core.MockDB) {},
			wordBehaviour: func(words *mock_core.MockWords) {},
			wantErr:       true,
		},
		{
			name:        "no matches",
			inputPhrase: "aboba",
			inputLimit:  1,
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().GetComicsIdsByWords(t.Context(), []string{"aboba"}).Return([]int{}, ErrNotFound)
			},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "aboba").Return([]string{"aboba"}, nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			db := mock_core.NewMockDB(c)
			words := mock_core.NewMockWords(c)

			testCase.wordBehaviour(words)
			testCase.dbBehaviour(db)

			service := NewService(slog.Default(), db, words, nil, nil, nil)

			res, err := service.Search(t.Context(), testCase.inputLimit, testCase.inputPhrase)

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.ElementsMatch(t, testCase.expected, res)
		})
	}
}

func TestIndexSearch(t *testing.T) {
	testCases := []struct {
		name           string
		inputPhrase    string
		inputLimit     int
		expected       []Comics
		dbBehaviour    func(db *mock_core.MockDB)
		wordBehaviour  func(words *mock_core.MockWords)
		indexBehaviour func(index *mock_core.MockIndex)
		wantErr        bool
	}{
		{
			name:        "OK",
			inputPhrase: "abcd",
			inputLimit:  1,
			expected:    []Comics{{1, ""}},
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().GetComicsesInfoByIds(t.Context(), []int{1}).Return([]ComicsInfo{{1, "", 4}}, nil)
			},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "abcd").Return([]string{"abcd"}, nil)
			},
			indexBehaviour: func(index *mock_core.MockIndex) {
				index.EXPECT().GetComicsIdsByWords(t.Context(), []string{"abcd"}).Return([]int{1}, nil)
			},
		},
		{
			name:        "test ranging",
			inputPhrase: "abcd",
			inputLimit:  1,
			expected:    []Comics{{1, ""}},
			dbBehaviour: func(db *mock_core.MockDB) {
				db.EXPECT().GetComicsesInfoByIds(t.Context(), []int{1, 2}).Return([]ComicsInfo{{1, "", 4}, {2, "", 8}}, nil)
			},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "abcd").Return([]string{"abcd"}, nil)
			},
			indexBehaviour: func(index *mock_core.MockIndex) {
				index.EXPECT().GetComicsIdsByWords(t.Context(), []string{"abcd"}).Return([]int{1, 2}, nil)
			},
		},
		{
			name:        "index error",
			inputPhrase: "abcd",
			inputLimit:  1,
			expected:    nil,
			dbBehaviour: func(db *mock_core.MockDB) {},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "abcd").Return([]string{"abcd"}, nil)
			},
			indexBehaviour: func(index *mock_core.MockIndex) {
				index.EXPECT().GetComicsIdsByWords(t.Context(), []string{"abcd"}).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name:        "missing phrase",
			inputLimit:  1,
			dbBehaviour: func(db *mock_core.MockDB) {},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "").Return(nil, errors.New("some errror"))
			},
			indexBehaviour: func(index *mock_core.MockIndex) {},
			wantErr:        true,
		},
		{
			name:           "missing limit",
			inputPhrase:    "aaaa",
			dbBehaviour:    func(db *mock_core.MockDB) {},
			wordBehaviour:  func(words *mock_core.MockWords) {},
			indexBehaviour: func(index *mock_core.MockIndex) {},
			wantErr:        true,
		},
		{
			name:        "no matches",
			inputPhrase: "aboba",
			inputLimit:  1,
			dbBehaviour: func(db *mock_core.MockDB) {},
			wordBehaviour: func(words *mock_core.MockWords) {
				words.EXPECT().Norm(t.Context(), "aboba").Return([]string{"aboba"}, nil)
			},
			indexBehaviour: func(index *mock_core.MockIndex) {
				index.EXPECT().GetComicsIdsByWords(t.Context(), []string{"aboba"}).Return([]int{}, ErrNotFound)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			db := mock_core.NewMockDB(c)
			words := mock_core.NewMockWords(c)
			index := mock_core.NewMockIndex(c)

			testCase.wordBehaviour(words)
			testCase.dbBehaviour(db)
			testCase.indexBehaviour(index)

			service := NewService(slog.Default(), db, words, index, nil, nil)

			res, err := service.SearchIndex(t.Context(), testCase.inputLimit, testCase.inputPhrase)

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.ElementsMatch(t, testCase.expected, res)
		})
	}
}
