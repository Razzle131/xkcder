package grpc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	searchpb "yadro.com/course/proto/search"
	"yadro.com/course/search/core"
	mock_core "yadro.com/course/search/core/mocks"
)

func TestSearch(t *testing.T) {
	testCases := []struct {
		name          string
		input         *searchpb.SearchRequest
		mockBehaviour func(m *mock_core.MockSearcher)
		expected      *searchpb.SearchReply
		wantErr       bool
	}{
		{
			name: "OK",
			input: &searchpb.SearchRequest{
				Limit:  1,
				Phrase: "abc",
			},
			mockBehaviour: func(m *mock_core.MockSearcher) {
				m.EXPECT().Search(t.Context(), 1, "abc").Return([]core.Comics{{ID: 1, URL: ""}}, nil)
			},
			expected: &searchpb.SearchReply{
				Comics: []*searchpb.Comics{{Id: 1, Url: ""}},
			},
		},
		{
			name: "error",
			input: &searchpb.SearchRequest{
				Limit:  1,
				Phrase: "abc",
			},
			mockBehaviour: func(m *mock_core.MockSearcher) {
				m.EXPECT().Search(t.Context(), 1, "abc").Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_core.NewMockSearcher(c)
			testCase.mockBehaviour(service)

			srv := NewServer(service)

			res, err := srv.Search(t.Context(), testCase.input)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, res)
		})
	}
}

func TestSearchIndex(t *testing.T) {
	testCases := []struct {
		name          string
		input         *searchpb.SearchRequest
		mockBehaviour func(m *mock_core.MockSearcher)
		expected      *searchpb.SearchReply
		wantErr       bool
	}{
		{
			name: "OK",
			input: &searchpb.SearchRequest{
				Limit:  1,
				Phrase: "abc",
			},
			mockBehaviour: func(m *mock_core.MockSearcher) {
				m.EXPECT().SearchIndex(t.Context(), 1, "abc").Return([]core.Comics{{ID: 1, URL: ""}}, nil)
			},
			expected: &searchpb.SearchReply{
				Comics: []*searchpb.Comics{{Id: 1, Url: ""}},
			},
		},
		{
			name: "error",
			input: &searchpb.SearchRequest{
				Limit:  1,
				Phrase: "abc",
			},
			mockBehaviour: func(m *mock_core.MockSearcher) {
				m.EXPECT().SearchIndex(t.Context(), 1, "abc").Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_core.NewMockSearcher(c)
			testCase.mockBehaviour(service)

			srv := NewServer(service)

			res, err := srv.SearchIndex(t.Context(), testCase.input)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, res)
		})
	}
}
