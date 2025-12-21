package aaa

import (
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt"

	"yadro.com/course/api/core"
)

const secretKey = "something secret here" // token sign key
const adminRole = "superuser"             // token subject

// Authentication, Authorization, Accounting
type AAA struct {
	users    map[string]string
	tokenTTL time.Duration
	log      *slog.Logger
}

func New(tokenTTL time.Duration, user, password string, log *slog.Logger) AAA {
	return AAA{
		users:    map[string]string{user: password},
		tokenTTL: tokenTTL,
		log:      log,
	}
}

type jwtClaims struct {
	Subject string
	jwt.StandardClaims
}

func (a AAA) Login(name, password string) (string, error) {
	userPassword, found := a.users[name]
	if !found {
		return "", core.ErrNotAuthorized
	}

	if userPassword != password {
		return "", core.ErrNotAuthorized
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		adminRole,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}).SignedString([]byte(secretKey))
	if err != nil {
		a.log.Error("failed to sign token", "err", err)
		return "", core.ErrNotAuthorized
	}

	return token, nil
}

func (a AAA) Verify(tokenString string) error {
	claims := jwtClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		slog.Error("failed token validation", "err", err)
		return core.ErrNotAuthorized
	}
	if !parsedToken.Valid {
		return core.ErrNotAuthorized
	}

	return nil
}
