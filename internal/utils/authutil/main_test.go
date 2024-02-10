package authutil

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"gotube/internal/config"
	"gotube/pkg/model"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateTokenForUser(t *testing.T) {
	user := model.User{
		ID:        1,
		Name:      "MOhammad",
		Email:     "t@t.com",
		Password:  "12345678",
		CreatedAt: time.Now(),
	}

	conf := config.Data{
		JWTSecret:        "secret",
		JWTExpireMinutes: 2,
		Domain:           "http://gotube.com",
	}

	token, err := CreateTokenForUser(&user, conf)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotZero(t, token.ExpiresAt)

}

func TestVerifyAuthTokenInRequestHeader(t *testing.T) {
	conf := config.Data{
		JWTSecret:        "secret",
		JWTExpireMinutes: 2,
		Domain:           "http://gotube.com",
	}

	tests := []struct {
		name           string
		token          string
		expectedErr    string
		authHeader     bool
		expectedUserID int
	}{
		{
			name:           "Valid Token",
			token:          generateTestToken(conf),
			authHeader:     true,
			expectedUserID: 1,
		},
		{
			name:        "Invalid Token",
			token:       "invalidtoken",
			expectedErr: InvalidTokenErr.Error(),
			authHeader:  true,
		},
		{
			name:        "Expired token",
			token:       generateExpiredTestToken(conf),
			expectedErr: ExpiredTokenErr.Error(),
			authHeader:  true,
		},
		{
			name:        "missing auth header",
			token:       generateExpiredTestToken(conf),
			expectedErr: InvalidAuthHeaderErr.Error(),
			authHeader:  false,
		},
		{
			name:        "missing auth header's token",
			token:       "",
			expectedErr: InvalidAuthHeaderErr.Error(),
			authHeader:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// create sample request
			req := httptest.NewRequest("GET", "/", nil)

			if tc.authHeader {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}

			// call the verifier
			claims, err := VerifyAuthTokenInRequestHeader(req, conf)

			// check for expected error
			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			if tc.expectedUserID != 0 {
				assert.Equal(t, tc.expectedUserID, claims.UserID)
			}
		})
	}
}

func generateTestToken(config config.Data) string {
	token, _ := CreateTokenForUser(&model.User{ID: 1}, config)
	return token.PlainText
}

func generateExpiredTestToken(config config.Data) string {
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Unix(),
	})
	signedExpiredToken, _ := expiredToken.SignedString([]byte(config.JWTSecret))
	return signedExpiredToken
}
