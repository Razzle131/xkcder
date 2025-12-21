package xkcd

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"yadro.com/course/update/core"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		name      string
		xkcdReply xkcdGetReply
		xkcdCode  int
		expected  core.XKCDInfo
		wantErr   bool
	}{
		{
			name: "OK",
			xkcdReply: xkcdGetReply{
				ID: 1,
			},
			xkcdCode: http.StatusOK,
			expected: core.XKCDInfo{
				ID: 1,
			},
		},
		{
			name:     "not found",
			xkcdCode: http.StatusNotFound,
			wantErr:  true,
		},
		{
			name:     "internal xkcd code",
			xkcdCode: http.StatusInternalServerError,
			wantErr:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.xkcdCode)
				err := json.NewEncoder(w).Encode(testCase.xkcdReply)
				assert.NoError(t, err)
			}))
			defer srv.Close()

			c, err := NewClient(srv.URL, time.Second, slog.Default())
			require.NoError(t, err)

			res, err := c.Get(t.Context(), 1)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, res)
		})
	}
}

func TestLastID(t *testing.T) {
	testCases := []struct {
		name      string
		xkcdReply xkcdLastIdReply
		xkcdCode  int
		expected  int
		wantErr   bool
	}{
		{
			name: "OK",
			xkcdReply: xkcdLastIdReply{
				LastId: 1,
			},
			xkcdCode: http.StatusOK,
			expected: 1,
		},
		{
			name:     "internal xkcd code",
			xkcdCode: http.StatusInternalServerError,
			wantErr:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.xkcdCode)
				err := json.NewEncoder(w).Encode(testCase.xkcdReply)
				assert.NoError(t, err)
			}))
			defer srv.Close()

			c, err := NewClient(srv.URL, time.Second, slog.Default())
			require.NoError(t, err)

			res, err := c.LastID(t.Context())
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, res)
		})
	}
}
