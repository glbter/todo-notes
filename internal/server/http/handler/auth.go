package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"todoNote/internal/model"
	in_memory "todoNote/internal/repo/in-memory"
	"todoNote/internal/server/http/dto"
	"todoNote/internal/server/http/log"
	"todoNote/internal/usecase"
)


type Auth struct {
	usecaseUser usecase.IUserUsecase
	auth IAuth
	log log.Logger
}

type IAuth interface {
	CreateToken(user model.UserInReq) (string, error)
	ValidateToken(t string) (model.UserInReq, error)
}

func NewAuthHandler(u usecase.IUserUsecase, a IAuth, log log.Logger) *Auth {
	return &Auth{
		usecaseUser: u,
		auth: a,
		log: log,
	}
}

func(h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var u dto.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeErrorMessage(w, http.StatusBadRequest, wrongBody)
		return
	}

	usr, err := h.usecaseUser.FindByName(r.Context(), u.UserName)
	if errors.As(err, &in_memory.NoSuchNameError{}) {
		writeErrorMessage(w, http.StatusBadRequest, incorrectLoginOrPassword)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(fmt.Sprintf("login: hash password: user(name: %v) error: %v", u.UserName, err))
		return
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(u.Password)); err != nil {
		writeErrorMessage(w, http.StatusBadRequest, incorrectLoginOrPassword)
		return
	}

	tk, err := h.auth.CreateToken(model.UserInReq{Id: usr.Id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(fmt.Sprintf("login: create token: user(id: %v) error: %v", usr.Id, err))
		return
	}

	json.NewEncoder(w).Encode(dto.NewTokenBearer(tk))
}
