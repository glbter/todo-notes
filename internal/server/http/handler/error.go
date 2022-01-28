package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"todoNote/internal/model"
	"todoNote/internal/server/http/dto"
)

const (
	wrongBody                = "bad body format"
	wrongPathParams = "bad path parameter"
	incorrectLoginOrPassword = "incorrect login or password"
	passwordNotEqual = "password and confirm password are not equal"
	wrongDateFormat = "invalid date-time format"

	noNoteFound = "no such note found"
)

func getIdFromRequest(r *http.Request, urlParam string) (model.Id, error){
	urlId := chi.URLParam(r, noteIdParam)
	id, err := strconv.Atoi(urlId)
	if err != nil || id < 0 {
		return 0, fmt.Errorf("bad id path param")
	}

	return model.Id(id), nil
}

func writeErrorMessage(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(
		dto.Error{
			Message: msg,
			Error: code,
		})

	return
}

