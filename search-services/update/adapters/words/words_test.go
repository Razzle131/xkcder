package words

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wordspb "yadro.com/course/proto/words"
	mock_words "yadro.com/course/proto/words/mocks"
)

func TestPing(t *testing.T) {
	testCases := []struct {
		name          string
		mockBehaviour func(client *mock_words.MockWordsClient)
		wantErr       bool
	}{
		{
			name: "OK",
			mockBehaviour: func(client *mock_words.MockWordsClient) {
				client.EXPECT().Ping(t.Context(), nil).Return(nil, nil)
			},
		},
		{
			name: "with error",
			mockBehaviour: func(client *mock_words.MockWordsClient) {
				client.EXPECT().Ping(t.Context(), nil).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			words := mock_words.NewMockWordsClient(c)
			testCase.mockBehaviour(words)

			client := Client{
				log:    slog.Default(),
				client: words,
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

func TestNorm(t *testing.T) {
	testCases := []struct {
		name          string
		inputPhrase   string
		expected      []string
		mockBehaviour func(client *mock_words.MockWordsClient)
		wantErr       bool
	}{
		{
			name:        "OK",
			inputPhrase: "abc",
			expected:    []string{"abc"},
			mockBehaviour: func(client *mock_words.MockWordsClient) {
				client.EXPECT().Norm(t.Context(), &wordspb.WordsRequest{
					Phrase: "abc",
				}).Return(&wordspb.WordsReply{Words: []string{"abc"}}, nil)
			},
		},
		{
			name:          "bad phrase",
			inputPhrase:   "",
			mockBehaviour: func(client *mock_words.MockWordsClient) {},
			wantErr:       true,
		},
		{
			name:        "with exausted error",
			inputPhrase: "abc",
			mockBehaviour: func(client *mock_words.MockWordsClient) {
				client.EXPECT().Norm(t.Context(), &wordspb.WordsRequest{
					Phrase: "abc",
				}).Return(nil, status.Error(codes.ResourceExhausted, ""))
			},
			wantErr: true,
		},
		{
			name:        "with random error",
			inputPhrase: "abc",
			mockBehaviour: func(client *mock_words.MockWordsClient) {
				client.EXPECT().Norm(t.Context(), &wordspb.WordsRequest{
					Phrase: "abc",
				}).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			words := mock_words.NewMockWordsClient(c)
			testCase.mockBehaviour(words)

			client := Client{
				log:    slog.Default(),
				client: words,
			}

			res, err := client.Norm(t.Context(), testCase.inputPhrase)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.ElementsMatch(t, testCase.expected, res)
		})
	}
}
