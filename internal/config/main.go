package config

import (
	"os"
	"strconv"
)

type Data struct {
	JWTSecret        string
	JWTExpireMinutes int
	Domain           string
	Port             string

	GoogleClientID     string
	GoogleClientSECRET string
	GoogleCallback     string
}

func New() Data {
	jwtExpire, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN_MINUTES"))
	return Data{
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpireMinutes:   jwtExpire,
		Domain:             os.Getenv("DOMAIN"),
		Port:               os.Getenv("PORT"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSECRET: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleCallback:     os.Getenv("GOOGLE_CALLBACK"),
	}
}
