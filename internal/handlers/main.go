package handlers

import (
	"gotube/internal/config"
	"gotube/internal/handlers/channelauth"
	"gotube/internal/handlers/userhandler"
	"gotube/pkg/repository"
)

type Handler struct {
	UsersHandler       userhandler.UsersHandler
	ChannelAuthHandler channelauth.Handler
	config             config.Data
}

func New(repo repository.Repository, config config.Data) Handler {
	return Handler{
		UsersHandler:       userhandler.NewHandler(repo.UserRepository, config),
		ChannelAuthHandler: channelauth.NewHandler(repo.ChannelRepository, config),
		config:             config,
	}
}
