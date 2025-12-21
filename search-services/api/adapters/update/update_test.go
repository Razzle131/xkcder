package update

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"yadro.com/course/api/core"
	updatepb "yadro.com/course/proto/update"
	mock_update "yadro.com/course/proto/update/mocks"
)

func TestPing(t *testing.T) {
	testCases := []struct {
		name          string
		mockBehaviour func(s *mock_update.MockUpdateClient)
		wantErr       bool
	}{
		{
			name: "OK",
			mockBehaviour: func(s *mock_update.MockUpdateClient) {
				s.EXPECT().Ping(t.Context(), nil).Return(nil, nil)
			},
		},
		{
			name: "with error",
			mockBehaviour: func(s *mock_update.MockUpdateClient) {
				s.EXPECT().Ping(t.Context(), nil).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			update := mock_update.NewMockUpdateClient(c)
			testCase.mockBehaviour(update)

			client := Client{
				log:    slog.Default(),
				client: update,
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

func TestStatus(t *testing.T) {
	testCases := []struct {
		name          string
		expected      core.UpdateStatus
		mockBehaviour func(s *mock_update.MockUpdateClient)
		wantErr       bool
	}{
		{
			name:     "OK",
			expected: core.StatusUpdateIdle,
			mockBehaviour: func(client *mock_update.MockUpdateClient) {
				client.EXPECT().Status(t.Context(), nil).Return(&updatepb.StatusReply{Status: updatepb.Status_STATUS_IDLE}, nil)
			},
		},
		{
			name:     "unknown status",
			expected: core.StatusUpdateUnknown,
			mockBehaviour: func(client *mock_update.MockUpdateClient) {
				client.EXPECT().Status(t.Context(), nil).Return(&updatepb.StatusReply{Status: updatepb.Status(228)}, nil)
			},
		},
		{
			name: "with error",
			mockBehaviour: func(client *mock_update.MockUpdateClient) {
				client.EXPECT().Status(t.Context(), nil).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			update := mock_update.NewMockUpdateClient(c)
			testCase.mockBehaviour(update)

			client := Client{
				log:    slog.Default(),
				client: update,
			}

			res, err := client.Status(t.Context())
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, res)
		})
	}
}

func TestStats(t *testing.T) {
	testCases := []struct {
		name          string
		expected      core.UpdateStats
		mockBehaviour func(s *mock_update.MockUpdateClient)
		wantErr       bool
	}{
		{
			name:     "OK",
			expected: core.UpdateStats{WordsTotal: 1, WordsUnique: 1, ComicsFetched: 1, ComicsTotal: 1},
			mockBehaviour: func(client *mock_update.MockUpdateClient) {
				client.EXPECT().Stats(t.Context(), nil).Return(&updatepb.StatsReply{WordsTotal: 1, WordsUnique: 1, ComicsFetched: 1, ComicsTotal: 1}, nil)
			},
		},
		{
			name: "with error",
			mockBehaviour: func(client *mock_update.MockUpdateClient) {
				client.EXPECT().Stats(t.Context(), nil).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			update := mock_update.NewMockUpdateClient(c)
			testCase.mockBehaviour(update)

			client := Client{
				log:    slog.Default(),
				client: update,
			}

			res, err := client.Stats(t.Context())
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, res)
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name          string
		mockBehaviour func(s *mock_update.MockUpdateClient)
		wantErr       bool
	}{
		{
			name: "OK",
			mockBehaviour: func(client *mock_update.MockUpdateClient) {
				client.EXPECT().Update(t.Context(), nil).Return(nil, nil)
			},
		},
		{
			name: "with random error",
			mockBehaviour: func(client *mock_update.MockUpdateClient) {
				client.EXPECT().Update(t.Context(), nil).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "with alredy exists error",
			mockBehaviour: func(client *mock_update.MockUpdateClient) {
				client.EXPECT().Update(t.Context(), nil).Return(nil, status.Error(codes.AlreadyExists, ""))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			update := mock_update.NewMockUpdateClient(c)
			testCase.mockBehaviour(update)

			client := Client{
				log:    slog.Default(),
				client: update,
			}

			err := client.Update(t.Context())
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDrop(t *testing.T) {
	testCases := []struct {
		name          string
		mockBehaviour func(s *mock_update.MockUpdateClient)
		wantErr       bool
	}{
		{
			name: "OK",
			mockBehaviour: func(client *mock_update.MockUpdateClient) {
				client.EXPECT().Drop(t.Context(), nil).Return(nil, nil)
			},
		},
		{
			name: "with error",
			mockBehaviour: func(client *mock_update.MockUpdateClient) {
				client.EXPECT().Drop(t.Context(), nil).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			update := mock_update.NewMockUpdateClient(c)
			testCase.mockBehaviour(update)

			client := Client{
				log:    slog.Default(),
				client: update,
			}

			err := client.Drop(t.Context())
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
