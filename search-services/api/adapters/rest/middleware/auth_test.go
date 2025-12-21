package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_middleware "yadro.com/course/api/adapters/rest/middleware/mocks"
)

func TestAuth(t *testing.T) {
	testCases := []struct {
		name          string
		token         string
		expectedCode  int
		mockBehaviour func(m *mock_middleware.MockTokenVerifier)
	}{
		{
			name:         "OK",
			token:        "Token abcd",
			expectedCode: http.StatusOK,
			mockBehaviour: func(m *mock_middleware.MockTokenVerifier) {
				m.EXPECT().Verify("abcd").Return(nil)
			},
		},
		{
			name:          "not match regexp",
			token:         "abcd",
			expectedCode:  http.StatusUnauthorized,
			mockBehaviour: func(m *mock_middleware.MockTokenVerifier) {},
		},
		{
			name:         "not verified",
			token:        "Token abcd",
			expectedCode: http.StatusUnauthorized,
			mockBehaviour: func(m *mock_middleware.MockTokenVerifier) {
				m.EXPECT().Verify("abcd").Return(errors.New("some error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest("GET", "/api/auth", nil)
			req.Header[tokenHeader] = []string{testCase.token}

			verifier := mock_middleware.NewMockTokenVerifier(c)
			testCase.mockBehaviour(verifier)

			w := httptest.NewRecorder()
			Auth(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}, verifier)(w, req)

			res := w.Result()
			assert.Equal(t, testCase.expectedCode, res.StatusCode)
		})
	}
}
