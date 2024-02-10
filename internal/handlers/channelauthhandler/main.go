package channelauthhandler

import (
	"encoding/json"
	"gotube/internal/config"
	"gotube/internal/handlers/handlerutils"
	"gotube/internal/utils/channelauth"
	"gotube/pkg/model"
	"gotube/pkg/repository"
	"net/http"
)

type Handler struct {
	repo   repository.ChannelRepository // channel repo
	config config.Data
	util   channelauth.Util
}

func New(
	repo repository.ChannelRepository,
	config config.Data,
	tokenExchanger channelauth.TokenExchanger,
	channelLister channelauth.ChannelLister,
) Handler {

	return Handler{
		repo:   repo,
		config: config,
		util:   channelauth.New(config, tokenExchanger, channelLister),
	}
}

type url struct {
	Plain string `json:"plain"`
}

func (h *Handler) GenerateAuthLink(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*model.User)
	if !ok {
		handlerutils.ReturnJsonError(w, http.StatusUnauthorized, "failed to get auth user")
		return
	}

	urlPlain := h.util.GenerateAuthUrl(*user)

	url := url{Plain: urlPlain}

	if err := json.NewEncoder(w).Encode(&url); err != nil {
		handlerutils.ReturnJsonError(w, http.StatusBadRequest, "failed to get link")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	channel, err := h.util.HandleCallback(r)
	if err != nil {
		handlerutils.ReturnJsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	insertedChannel, err := h.repo.UpdateOrCreate(r.Context(), *channel)
	if err != nil {
		handlerutils.ReturnJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handlerutils.ReturnJson(w, 200, insertedChannel)
}
