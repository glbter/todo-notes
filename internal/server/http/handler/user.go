package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"todoNote/internal/model"
	"todoNote/internal/repo/postgres"
	"todoNote/internal/server/http/dto"
	"todoNote/internal/server/http/log"
	"todoNote/internal/server/http/middleware"
	"todoNote/internal/usecase"
)


type User struct {
	userCase usecase.IUserUsecase
	log log.Logger
}

func NewUserHandler(u usecase.IUserUsecase, log log.Logger) *User {
	return &User{
		userCase: u,
		log: log,
	}
}

func(h *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u dto.UserRegistration
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeErrorMessage(w, http.StatusBadRequest, wrongBody)
		return
	}

	if _, ok:= usecase.ValidateZone(u.TimeZone); !ok {
		writeErrorMessage(w, http.StatusBadRequest, wrongDateFormat)
		return
	}

	id, err := h.userCase.Create(r.Context(), &model.UserNew{
		Name: u.UserName,
		TimeZone: u.TimeZone,
		Password: u.Password,
	})

	if errors.As(err, &postgres.UserExistsError{}) {
		writeErrorMessage(w, http.StatusBadRequest, "such username is taken")
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error("create user: db err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.IdObject{Id: id})
}

func(h *User) PartialUpdateUser(w http.ResponseWriter, r *http.Request) {
	var u dto.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeErrorMessage(w, http.StatusBadRequest, wrongBody)
		return
	}

	if _, ok:= usecase.ValidateZone(u.TimeZone); !ok {
		writeErrorMessage(w, http.StatusBadRequest, wrongDateFormat)
		return
	}

	usr, ok := middleware.UserFromContext(r, h.log,"update user")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := model.UserUpdate{Id: usr.Id, TimeZone: u.TimeZone}
	if err := h.userCase.Update(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error("patch update user: user(id: %v) db err: %v", usr.Id, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func(h *User) DeleteUser(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.UserFromContext(r, h.log,"delete user")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.userCase.Remove(r.Context(), u.Id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error("remove user: user(id: %v) db err: %v", u.Id, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


