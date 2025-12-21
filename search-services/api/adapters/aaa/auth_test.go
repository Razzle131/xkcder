package aaa

import (
	"log/slog"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const adminUser = "ADMIN_USER"
const adminPass = "ADMIN_PASSWORD"

func TestLogin(t *testing.T) {
	testCases := []struct {
		name          string
		inputName     string
		inputPassword string
		wantErr       bool
	}{
		{
			name:          "OK",
			inputName:     adminUser,
			inputPassword: adminPass,
		},
		{
			name:      "missing password",
			inputName: adminUser,
			wantErr:   true,
		},
		{
			name:          "missing username",
			inputPassword: adminPass,
			wantErr:       true,
		},
		{
			name:          "user is not registered",
			inputName:     "foo",
			inputPassword: "foo",
			wantErr:       true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			aaa := New(time.Hour, adminUser, adminPass, slog.Default())

			_, err := aaa.Login(testCase.inputName, testCase.inputPassword)

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestVerify(t *testing.T) {
	testCases := []struct {
		name          string
		inputName     string
		inputPassword string
		waitTime      time.Duration
		wantErr       bool
	}{
		{
			name:          "OK",
			inputName:     adminUser,
			inputPassword: adminPass,
		},
		{
			name:          "token expired recently",
			inputName:     adminUser,
			inputPassword: adminPass,
			waitTime:      2 * time.Second,
			wantErr:       true,
		},
		{
			name:          "old token",
			inputName:     adminUser,
			inputPassword: adminPass,
			waitTime:      5 * time.Hour,
			wantErr:       true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				aaa := New(time.Second, adminUser, adminPass, slog.Default())

				token, err := aaa.Login(testCase.inputName, testCase.inputPassword)
				require.NoError(t, err)

				time.Sleep(testCase.waitTime)

				err = aaa.Verify(token)
				if testCase.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
			})

		})
	}
}
