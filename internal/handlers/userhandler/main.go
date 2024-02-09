package userhandler

import (
	"gotube/internal/config"
	"gotube/internal/handlers/handlerutils"
	"gotube/internal/utils/authutil"
	"gotube/internal/utils/passwordutil"
	"gotube/pkg/model"
	"gotube/pkg/repository"
	"net/http"
)

type UsersHandler struct {
	repo   repository.UserRepository
	config config.Data
}

func NewHandler(repo repository.UserRepository, config config.Data) UsersHandler {
	return UsersHandler{
		repo:   repo,
		config: config,
	}
}

// Register registers a new user
func (h *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		handlerutils.ReturnJsonError(w, http.StatusBadRequest, "Email is invalid")
		return
	}

	name := r.FormValue("name")
	if name == "" {
		handlerutils.ReturnJsonError(w, http.StatusBadRequest, "Name is invalid")
		return
	}

	password := r.FormValue("password")
	if password == "" {
		handlerutils.ReturnJsonError(w, http.StatusBadRequest, "Password is invalid")
		return
	}

	// check to see email not exists
	exists, err := h.repo.EmailExists(email)
	if err != nil {
		handlerutils.ReturnJsonError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	if exists {
		handlerutils.ReturnJsonError(w, http.StatusBadRequest, "Email already exists")
		return
	}

	// get name, email and password from request
	// and create a new user
	err = h.repo.Create(r.Context(), model.User{
		Name:     name,
		Email:    email,
		Password: password,
	})

	if err != nil {
		handlerutils.ReturnJsonError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	handlerutils.ReturnJsonMessages(w, http.StatusOK, "User created successfully")
}

func (h *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		handlerutils.ReturnJsonError(w, http.StatusBadRequest, "Email is invalid")
		return
	}

	password := r.FormValue("password")
	if password == "" {
		handlerutils.ReturnJsonError(w, http.StatusBadRequest, "Password is invalid")
		return
	}

	// get user from repo by email
	user, err := h.repo.FindByField(r.Context(), "email", email)
	if err != nil {
		handlerutils.ReturnJsonError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	if user == nil {
		handlerutils.ReturnJsonError(w, http.StatusNotFound, "Email or Password is invalid")
		return
	}
	// compare password and hash
	passwordMatches := passwordutil.PasswordMatches(user.Password, password)

	if !passwordMatches {
		handlerutils.ReturnJsonError(w, http.StatusUnauthorized, "Email or Password is invalid")
		return
	}
	// create access Token
	token, err := authutil.CreateTokenForUser(user, h.config)
	if err != nil {
		handlerutils.ReturnJsonError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	handlerutils.ReturnJson(w, http.StatusOK, token)
}
