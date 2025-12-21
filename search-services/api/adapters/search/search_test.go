package search

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"yadro.com/course/api/core"
	searchpb "yadro.com/course/proto/search"
	mock_search "yadro.com/course/proto/search/mocks"
)

func TestPing(t *testing.T) {
	testCases := []struct {
		name          string
		mockBehaviour func(s *mock_search.MockSearchClient)
		wantErr       bool
	}{
		{
			name: "OK",
			mockBehaviour: func(s *mock_search.MockSearchClient) {
				s.EXPECT().Ping(t.Context(), nil).Return(nil, nil)
			},
		},
		{
			name: "with error",
			mockBehaviour: func(s *mock_search.MockSearchClient) {
				s.EXPECT().Ping(t.Context(), nil).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			search := mock_search.NewMockSearchClient(c)
			testCase.mockBehaviour(search)

			client := Client{
				log:    slog.Default(),
				client: search,
			}

			err := client.Ping(t.Context())
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestSearch(t *testing.T) {
	testCases := []struct {
		name          string
		inputPhrase   string
		inputLimit    int
		expected      []core.Comics
		mockBehaviour func(s *mock_search.MockSearchClient)
		wantErr       bool
	}{
		{
			name:        "OK",
			inputPhrase: "abc",
			inputLimit:  1,
			expected:    []core.Comics{{ID: 1, URL: ""}},
			mockBehaviour: func(s *mock_search.MockSearchClient) {
				s.EXPECT().Search(t.Context(), &searchpb.SearchRequest{Limit: 1, Phrase: "abc"}).Return(&searchpb.SearchReply{Comics: []*searchpb.Comics{&searchpb.Comics{Id: 1, Url: ""}}}, nil)
			},
		},
		{
			name:        "not found",
			inputPhrase: "abc",
			inputLimit:  1,
			expected:    []core.Comics{},
			mockBehaviour: func(s *mock_search.MockSearchClient) {
				s.EXPECT().Search(t.Context(), &searchpb.SearchRequest{Limit: 1, Phrase: "abc"}).Return(&searchpb.SearchReply{Comics: []*searchpb.Comics{}}, nil)
			},
		},
		{
			name:        "with error",
			inputPhrase: "abc",
			inputLimit:  1,
			mockBehaviour: func(s *mock_search.MockSearchClient) {
				s.EXPECT().Search(t.Context(), &searchpb.SearchRequest{Limit: 1, Phrase: "abc"}).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			search := mock_search.NewMockSearchClient(c)
			testCase.mockBehaviour(search)

			client := Client{
				log:    slog.Default(),
				client: search,
			}

			comicses, err := client.Search(t.Context(), testCase.inputPhrase, testCase.inputLimit)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.ElementsMatch(t, testCase.expected, comicses)
		})
	}
}

func TestSearchIndex(t *testing.T) {
	testCases := []struct {
		name          string
		inputPhrase   string
		inputLimit    int
		expected      []core.Comics
		mockBehaviour func(s *mock_search.MockSearchClient)
		wantErr       bool
	}{
		{
			name:        "OK",
			inputPhrase: "abc",
			inputLimit:  1,
			expected:    []core.Comics{{ID: 1, URL: ""}},
			mockBehaviour: func(s *mock_search.MockSearchClient) {
				s.EXPECT().SearchIndex(t.Context(), &searchpb.SearchRequest{Limit: 1, Phrase: "abc"}).Return(&searchpb.SearchReply{Comics: []*searchpb.Comics{&searchpb.Comics{Id: 1, Url: ""}}}, nil)
			},
		},
		{
			name:        "not found",
			inputPhrase: "abc",
			inputLimit:  1,
			expected:    []core.Comics{},
			mockBehaviour: func(s *mock_search.MockSearchClient) {
				s.EXPECT().SearchIndex(t.Context(), &searchpb.SearchRequest{Limit: 1, Phrase: "abc"}).Return(&searchpb.SearchReply{Comics: []*searchpb.Comics{}}, nil)
			},
		},
		{
			name:        "with error",
			inputPhrase: "abc",
			inputLimit:  1,
			mockBehaviour: func(s *mock_search.MockSearchClient) {
				s.EXPECT().SearchIndex(t.Context(), &searchpb.SearchRequest{Limit: 1, Phrase: "abc"}).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			search := mock_search.NewMockSearchClient(c)
			testCase.mockBehaviour(search)

			client := Client{
				log:    slog.Default(),
				client: search,
			}

			comicses, err := client.SearchIndex(t.Context(), testCase.inputPhrase, testCase.inputLimit)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.ElementsMatch(t, testCase.expected, comicses)
		})
	}
}
