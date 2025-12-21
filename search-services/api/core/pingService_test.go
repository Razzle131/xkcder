package core_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	. "yadro.com/course/api/core"
	mock_core "yadro.com/course/api/core/mocks"
)

func TestPingServices(t *testing.T) {
	testCases := []struct {
		name         string
		prepareInput func(c *gomock.Controller) map[string]Pinger
		expected     map[string]string
	}{
		{
			name: "OK",
			prepareInput: func(c *gomock.Controller) map[string]Pinger {
				pingers := make(map[string]Pinger)
				m1 := mock_core.NewMockPinger(c)
				m1.EXPECT().Ping(t.Context()).Return(nil)
				pingers["a"] = m1
				return pingers
			},
			expected: map[string]string{
				"a": ServiceStatusOk,
			},
		},
		{
			name: "empty",
			prepareInput: func(c *gomock.Controller) map[string]Pinger {
				return map[string]Pinger{}
			},
			expected: map[string]string{},
		},
		{
			name: "only unavailable",
			prepareInput: func(c *gomock.Controller) map[string]Pinger {
				pingers := make(map[string]Pinger)
				m1 := mock_core.NewMockPinger(c)
				m1.EXPECT().Ping(t.Context()).Return(errors.New("some error"))
				pingers["a"] = m1
				return pingers
			},
			expected: map[string]string{
				"a": ServiceStatusUnavailable,
			},
		},
		{
			name: "composite",
			prepareInput: func(c *gomock.Controller) map[string]Pinger {
				pingers := make(map[string]Pinger)
				m1 := mock_core.NewMockPinger(c)
				m1.EXPECT().Ping(t.Context()).Return(nil)
				pingers["a"] = m1

				m2 := mock_core.NewMockPinger(c)
				m2.EXPECT().Ping(t.Context()).Return(errors.New("some error"))
				pingers["b"] = m2

				return pingers
			},
			expected: map[string]string{
				"a": ServiceStatusOk,
				"b": ServiceStatusUnavailable,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			pingers := testCase.prepareInput(c)
			res := PingServices(t.Context(), pingers)

			assert.Equal(t, testCase.expected, res)
		})
	}
}
