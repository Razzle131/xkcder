package grpc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	updatepb "yadro.com/course/proto/update"
	"yadro.com/course/update/core"
	mock_core "yadro.com/course/update/core/mocks"
)

func TestStatus(t *testing.T) {
	testCases := []struct {
		name          string
		mockBehaviour func(m *mock_core.MockUpdater)
		expected      *updatepb.StatusReply
	}{
		{
			name: "OK",
			mockBehaviour: func(m *mock_core.MockUpdater) {
				m.EXPECT().Status(t.Context()).Return(core.StatusIdle)
			},
			expected: &updatepb.StatusReply{
				Status: updatepb.Status_STATUS_IDLE,
			},
		},
		{
			name: "random status",
			mockBehaviour: func(m *mock_core.MockUpdater) {
				m.EXPECT().Status(t.Context()).Return(core.ServiceStatus("abc"))
			},
			expected: &updatepb.StatusReply{
				Status: updatepb.Status_STATUS_UNSPECIFIED,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_core.NewMockUpdater(c)
			testCase.mockBehaviour(service)

			srv := NewServer(service)

			res, err := srv.Status(t.Context(), nil)
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, res)
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name          string
		mockBehaviour func(m *mock_core.MockUpdater)
		wantErr       bool
	}{
		{
			name: "OK",
			mockBehaviour: func(m *mock_core.MockUpdater) {
				m.EXPECT().Update(t.Context()).Return(nil)
			},
		},
		{
			name: "error",
			mockBehaviour: func(m *mock_core.MockUpdater) {
				m.EXPECT().Update(t.Context()).Return(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_core.NewMockUpdater(c)
			testCase.mockBehaviour(service)

			srv := NewServer(service)

			_, err := srv.Update(t.Context(), nil)
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
		expected      *updatepb.StatsReply
		mockBehaviour func(m *mock_core.MockUpdater)
		wantErr       bool
	}{
		{
			name: "OK",
			expected: &updatepb.StatsReply{
				WordsTotal:    1,
				WordsUnique:   1,
				ComicsFetched: 1,
				ComicsTotal:   1,
			},
			mockBehaviour: func(m *mock_core.MockUpdater) {
				m.EXPECT().Stats(t.Context()).Return(core.ServiceStats{
					DBStats: core.DBStats{
						WordsTotal:    1,
						WordsUnique:   1,
						ComicsFetched: 1,
					},
					ComicsTotal: 1,
				}, nil)
			},
		},
		{
			name: "error",
			mockBehaviour: func(m *mock_core.MockUpdater) {
				m.EXPECT().Stats(t.Context()).Return(core.ServiceStats{}, errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_core.NewMockUpdater(c)
			testCase.mockBehaviour(service)

			srv := NewServer(service)

			res, err := srv.Stats(t.Context(), nil)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, res)
		})
	}
}
