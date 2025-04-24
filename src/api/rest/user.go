package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"nullableocean-postupashki/src/api/rest/types"
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/usecases"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService usecases.User
}

func NewUserHandler(userService usecases.User) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(router *chi.Mux) {
	router.Post("/user/register", h.Register)
	router.Post("/user/login", h.Login)
}

// @Summary Register user
// @Description "Register user"
// @Tags User
// @Accept json
// @Produce json
// @Param data body types.RegisterRequestBody true "user register data"
// @Success 201
// @Failure default {object} types.ErrorResponse "error response"
// @Router /user/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	registerData, err := types.ExtractRegisterBody(r)
	if err != nil {
		types.ProcessError(w, err, http.StatusBadRequest)
		return
	}

	user := domain.User{
		Login:    registerData.Login,
		Password: registerData.Password,
	}

	_, err = h.userService.Create(user)
	if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, usecases.ErrUserExist) {
			code = http.StatusBadRequest
		}

		types.ProcessError(w, err, code)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Login Login
// @Description "Login user"
// @Tags User
// @Produce json
// @Param data body types.LoginRequestBody true "user login data"
// @Success 201 {object} types.LoginResponse
// @Failure default {object} types.ErrorResponse "error response"
// @Router /user/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginData, err := types.ExtractLoginBody(r)
	if err != nil {
		types.ProcessError(w, err, http.StatusBadRequest)
		return
	}

	creds := &domain.User{
		Login:    loginData.Login,
		Password: loginData.Password,
	}
	session, err := h.userService.SignIn(creds)
	if err != nil {
		types.ProcessError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&types.LoginResponse{Token: session.SessionId})
}
