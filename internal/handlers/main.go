package handlers

import (
	"gotube/internal/config"
	"gotube/internal/handlers/channelauthhandler"
	"gotube/internal/handlers/userhandler"
	"gotube/pkg/repository"
)

type Handler struct {
	UsersHandler       userhandler.Handler
	ChannelAuthHandler channelauthhandler.Handler
	config             config.Data
}

func New(repo repository.Repository, config config.Data) Handler {
	return Handler{
		UsersHandler:       userhandler.New(repo.UserRepository, config),
		ChannelAuthHandler: channelauthhandler.New(repo.ChannelRepository, config),
		config:             config,
	}
}
