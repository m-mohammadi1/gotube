package authutil

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gotube/internal/config"
	"gotube/pkg/model"
	"net/http"
	"strings"
	"time"
)

type Token struct {
	PlainText string    `json:"plain_text"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	InvalidAuthHeaderErr = errors.New("invalid auth header")
	ExpiredTokenErr      = errors.New("expired token")
	InvalidTokenErr      = errors.New("invalid token")
)

// CreateTokenForUser Create a Token with an expire_at time
func CreateTokenForUser(user *model.User, config config.Data) (Token, error) {
	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = config.Domain
	claims["iss"] = config.Domain

	// set expire time
	expireTime := time.Now().Add(time.Minute * time.Duration(config.JWTExpireMinutes))
	claims["exp"] = expireTime.Unix()

	// sign the token
	signedToken, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return Token{}, err
	}

	return Token{
		PlainText: signedToken,
		ExpiresAt: expireTime,
	}, nil
}

func VerifyAuthTokenInRequestHeader(r *http.Request, config config.Data) (*Claims, error) {
	// get auth header
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))

	// sanity check
	if authHeader == "" {
		return nil, InvalidAuthHeaderErr
	}

	// check token format
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return nil, InvalidAuthHeaderErr
	}

	// check headerParts first to have Bearer
	if headerParts[0] != "Bearer" {
		return nil, InvalidAuthHeaderErr
	}

	tokenString := headerParts[1]

	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// validate signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid sigin method: %v", token.Header["alg"])
		}

		return []byte(config.JWTSecret), nil
	})

	// check for token errors
	if err != nil {
		// catch expired tokens
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return nil, ExpiredTokenErr
		}

		return nil, InvalidTokenErr
	}

	// check the issuer to be correct
	if claims.Issuer != config.Domain {
		return nil, InvalidTokenErr
	}

	return claims, err
}
