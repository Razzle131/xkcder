package rest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"yadro.com/course/api/core"
	mock_core "yadro.com/course/api/core/mocks"
)

func TestPing(t *testing.T) {
	testCases := []struct {
		name         string
		expectedCode int
		expectedBody string
		prepareMock  func(ctx context.Context, c *gomock.Controller) map[string]core.Pinger
	}{
		{
			name:         "OK",
			expectedCode: http.StatusOK,
			expectedBody: fmt.Sprintf(`{"replies":{"a":"%s"}}`+"\n", core.ServiceStatusOk),
			prepareMock: func(ctx context.Context, c *gomock.Controller) map[string]core.Pinger {
				pinger := mock_core.NewMockPinger(c)
				pinger.EXPECT().Ping(ctx).Return(nil)

				return map[string]core.Pinger{"a": pinger}
			},
		},
		{
			name:         "only failures",
			expectedCode: http.StatusOK,
			expectedBody: fmt.Sprintf(`{"replies":{"a":"%s"}}`+"\n", core.ServiceStatusUnavailable),
			prepareMock: func(ctx context.Context, c *gomock.Controller) map[string]core.Pinger {
				pinger := mock_core.NewMockPinger(c)
				pinger.EXPECT().Ping(ctx).Return(errors.New("some errror"))

				return map[string]core.Pinger{"a": pinger}
			},
		},
		{
			name:         "empty",
			expectedCode: http.StatusOK,
			expectedBody: `{"replies":{}}` + "\n",
			prepareMock: func(ctx context.Context, c *gomock.Controller) map[string]core.Pinger {
				return map[string]core.Pinger{}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest("GET", "/api/ping", nil)

			w := httptest.NewRecorder()
			NewPingHandler(slog.Default(), testCase.prepareMock(req.Context(), c))(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, testCase.expectedCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedBody, string(data))
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockBehaviour func(ctx context.Context, m *mock_core.MockUpdater)
	}{
		{
			name:         "OK",
			expectedCode: http.StatusOK,
			expectedBody: "",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockUpdater) {
				m.EXPECT().Update(ctx).Return(nil)
			},
		},
		{
			name:         "already updating",
			expectedCode: http.StatusAccepted,
			expectedBody: "",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockUpdater) {
				m.EXPECT().Update(ctx).Return(core.ErrAlreadyExists)
			},
		},
		{
			name:         "random error",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "updater error: some error\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockUpdater) {
				m.EXPECT().Update(ctx).Return(errors.New("some error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest("POST", "/api/db/update", nil)

			updater := mock_core.NewMockUpdater(c)
			testCase.mockBehaviour(req.Context(), updater)

			w := httptest.NewRecorder()
			NewUpdateHandler(slog.Default(), updater)(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, testCase.expectedCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedBody, string(data))
		})
	}
}

func TestUpdateStats(t *testing.T) {
	testCases := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockBehaviour func(ctx context.Context, m *mock_core.MockUpdater)
	}{
		{
			name:         "OK",
			expectedCode: http.StatusOK,
			expectedBody: `{"words_total":1,"words_unique":1,"comics_fetched":1,"comics_total":1}` + "\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockUpdater) {
				m.EXPECT().Stats(ctx).Return(core.UpdateStats{
					WordsTotal:    1,
					WordsUnique:   1,
					ComicsFetched: 1,
					ComicsTotal:   1,
				}, nil)
			},
		},
		{
			name:         "random error",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "updater error: some error\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockUpdater) {
				m.EXPECT().Stats(ctx).Return(core.UpdateStats{}, errors.New("some error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest("GET", "/api/db/stats", nil)

			updater := mock_core.NewMockUpdater(c)
			testCase.mockBehaviour(req.Context(), updater)

			w := httptest.NewRecorder()
			NewUpdateStatsHandler(slog.Default(), updater)(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, testCase.expectedCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedBody, string(data))
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	testCases := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockBehaviour func(ctx context.Context, m *mock_core.MockUpdater)
	}{
		{
			name:         "OK",
			expectedCode: http.StatusOK,
			expectedBody: fmt.Sprintf(`{"status":"%s"}`+"\n", core.StatusUpdateIdle),
			mockBehaviour: func(ctx context.Context, m *mock_core.MockUpdater) {
				m.EXPECT().Status(ctx).Return(core.StatusUpdateIdle, nil)
			},
		},
		{
			name:         "random error",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "updater error: some error\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockUpdater) {
				m.EXPECT().Status(ctx).Return(core.UpdateStatus(""), errors.New("some error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest("GET", "/api/db/status", nil)

			updater := mock_core.NewMockUpdater(c)
			testCase.mockBehaviour(req.Context(), updater)

			w := httptest.NewRecorder()
			NewUpdateStatusHandler(slog.Default(), updater)(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, testCase.expectedCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedBody, string(data))
		})
	}
}

func TestDrop(t *testing.T) {
	testCases := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockBehaviour func(ctx context.Context, m *mock_core.MockUpdater)
	}{
		{
			name:         "OK",
			expectedCode: http.StatusOK,
			expectedBody: "",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockUpdater) {
				m.EXPECT().Drop(ctx).Return(nil)
			},
		},
		{
			name:         "random error",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "updater error: some error\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockUpdater) {
				m.EXPECT().Drop(ctx).Return(errors.New("some error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest("DELETE", "/api/db", nil)

			updater := mock_core.NewMockUpdater(c)
			testCase.mockBehaviour(req.Context(), updater)

			w := httptest.NewRecorder()
			NewDropHandler(slog.Default(), updater)(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, testCase.expectedCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedBody, string(data))
		})
	}
}

func TestSearch(t *testing.T) {
	testCases := []struct {
		name          string
		inputLimit    int
		inputPhrase   string
		expectedCode  int
		expectedBody  string
		mockBehaviour func(ctx context.Context, m *mock_core.MockSearcher)
	}{
		{
			name:         "OK",
			inputLimit:   1,
			inputPhrase:  "abc",
			expectedCode: http.StatusOK,
			expectedBody: `{"comics":[{"id":1,"url":""}],"total":1}` + "\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockSearcher) {
				m.EXPECT().Search(ctx, "abc", 1).Return([]core.Comics{{ID: 1, URL: ""}}, nil)
			},
		},
		{
			name:         "no comicses found",
			inputLimit:   1,
			inputPhrase:  "abc",
			expectedCode: http.StatusOK,
			expectedBody: `{"comics":[],"total":0}` + "\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockSearcher) {
				m.EXPECT().Search(ctx, "abc", 1).Return([]core.Comics{}, nil)
			},
		},
		{
			name:          "empty limit",
			inputPhrase:   "abc",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "bad limit parameter\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockSearcher) {},
		},
		{
			name:          "empty phrase",
			inputLimit:    1,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "bad phrase parameter\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockSearcher) {},
		},
		{
			name:         "random error",
			inputLimit:   1,
			inputPhrase:  "abc",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "some error\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockSearcher) {
				m.EXPECT().Search(ctx, "abc", 1).Return(nil, errors.New("some error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/api/search?%s=%v&%s=%s", queryLimitParamName, testCase.inputLimit, queryPhraseParamName, testCase.inputPhrase),
				nil,
			)

			searcher := mock_core.NewMockSearcher(c)
			testCase.mockBehaviour(req.Context(), searcher)

			w := httptest.NewRecorder()
			NewSearchHandler(slog.Default(), searcher)(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, testCase.expectedCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedBody, string(data))
		})
	}
}

func TestIndexSearch(t *testing.T) {
	testCases := []struct {
		name          string
		inputLimit    int
		inputPhrase   string
		expectedCode  int
		expectedBody  string
		mockBehaviour func(ctx context.Context, m *mock_core.MockSearcher)
	}{
		{
			name:         "OK",
			inputLimit:   1,
			inputPhrase:  "abc",
			expectedCode: http.StatusOK,
			expectedBody: `{"comics":[{"id":1,"url":""}],"total":1}` + "\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockSearcher) {
				m.EXPECT().SearchIndex(ctx, "abc", 1).Return([]core.Comics{{ID: 1, URL: ""}}, nil)
			},
		},
		{
			name:         "no comicses found",
			inputLimit:   1,
			inputPhrase:  "abc",
			expectedCode: http.StatusOK,
			expectedBody: `{"comics":[],"total":0}` + "\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockSearcher) {
				m.EXPECT().SearchIndex(ctx, "abc", 1).Return([]core.Comics{}, nil)
			},
		},
		{
			name:          "empty limit",
			inputPhrase:   "abc",
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "bad limit parameter\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockSearcher) {},
		},
		{
			name:          "empty phrase",
			inputLimit:    1,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "bad phrase parameter\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockSearcher) {},
		},
		{
			name:         "random error",
			inputLimit:   1,
			inputPhrase:  "abc",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "some error\n",
			mockBehaviour: func(ctx context.Context, m *mock_core.MockSearcher) {
				m.EXPECT().SearchIndex(ctx, "abc", 1).Return(nil, errors.New("some error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/api/isearch?%s=%v&%s=%s", queryLimitParamName, testCase.inputLimit, queryPhraseParamName, testCase.inputPhrase),
				nil,
			)

			searcher := mock_core.NewMockSearcher(c)
			testCase.mockBehaviour(req.Context(), searcher)

			w := httptest.NewRecorder()
			NewIndexSearchHandler(slog.Default(), searcher)(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, testCase.expectedCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedBody, string(data))
		})
	}
}

func TestLogin(t *testing.T) {
	testCases := []struct {
		name          string
		inputBody     string
		expectedCode  int
		mockBehaviour func(ctx context.Context, m *mock_core.MockAAA)
		wantErr       bool
	}{
		{
			name:         "OK",
			inputBody:    `{"name":"admin", "password": "admin"}`,
			expectedCode: http.StatusOK,
			mockBehaviour: func(ctx context.Context, m *mock_core.MockAAA) {
				m.EXPECT().Login("admin", "admin").Return("token", nil)
			},
		},
		{
			name:         "unauthorized error",
			inputBody:    `{"name":"admin", "password": "admin"}`,
			expectedCode: http.StatusUnauthorized,
			mockBehaviour: func(ctx context.Context, m *mock_core.MockAAA) {
				m.EXPECT().Login("admin", "admin").Return("", core.ErrNotAuthorized)
			},
		},
		{
			name:         "random error",
			inputBody:    `{"name":"admin", "password": "admin"}`,
			expectedCode: http.StatusInternalServerError,
			mockBehaviour: func(ctx context.Context, m *mock_core.MockAAA) {
				m.EXPECT().Login("admin", "admin").Return("", errors.New("some error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer([]byte(testCase.inputBody)))

			aaa := mock_core.NewMockAAA(c)
			testCase.mockBehaviour(req.Context(), aaa)

			w := httptest.NewRecorder()
			NewLoginHandler(slog.Default(), aaa)(w, req)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, testCase.expectedCode, res.StatusCode)
		})
	}
}
