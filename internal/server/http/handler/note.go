package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todoNote/internal/model"
	"todoNote/internal/repo"
	"todoNote/internal/server/http/dto"
	"todoNote/internal/server/http/log"
	"todoNote/internal/server/http/middleware"
	"todoNote/internal/usecase"
)

const (
	noteIdParam = "noteId"

	timezoneQueryParam = "timezone"
	startFromQueryParam = "start_from"
	limitQueryParam = "limit"
	offsetQueryParam = "offset"
	isFinishedQueryParam = "is_finished"
)


type Note struct {
	usecaseNote usecase.INoteUsecase
	usecaseUser usecase.IUserUsecase
	log log.Logger
}

func NewNoteHandler(n usecase.INoteUsecase, u usecase.IUserUsecase, log log.Logger) *Note {
	return &Note{
		usecaseUser: u,
		usecaseNote: n,
		log: log,
	}
}


func(h *Note) CreateNote(w http.ResponseWriter, r *http.Request) {
	var n dto.NewNote
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeErrorMessage(w, http.StatusBadRequest, wrongBody)
		return
	}

	u, ok := middleware.UserFromContext(r, h.log, "create note")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error("create note: no user in context")
		return
	}

	note := model.NewNote(0, u.Id, n.Title, n.Text, n.Date, false)
	uId, err := h.usecaseNote.CreateNote(r.Context(), note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(fmt.Sprintf("create note: n: db err %v", err))
		return
	}

	json.NewEncoder(w).Encode(dto.IdObject{Id: uId})
	w.WriteHeader(http.StatusCreated)
}

func(h *Note) GetNotes(w http.ResponseWriter, r *http.Request) {
	startFrom := r.URL.Query().Get(startFromQueryParam)
	limit := r.URL.Query().Get(limitQueryParam)
	offset := r.URL.Query().Get(offsetQueryParam)
	timezone := r.URL.Query().Get(timezoneQueryParam)
	isFinished := r.URL.Query().Get(isFinishedQueryParam)

	u, ok := middleware.UserFromContext(r, h.log, "get notes")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p := usecase.FindParams{}
	zone, ok := h.checkZoneRule(r.Context(), timezone, u.Id)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p.Zone = zone
	p.Filter = repo.NoteFilter{}
	page := repo.PageFilter{
		Limit: repo.GetUIntParamPointer(offset),
		Offset: repo.GetUIntParamPointer(limit)}
	p.Filter.Page = page
	p.Filter.UserId = u.Id

	if startFrom != "" {
		if time, err := time.Parse(time.RFC3339, startFrom); err == nil {
			p.Filter.TakeFrom = &time
		}
	}
	if isFinished != "" {
		if b, err := strconv.ParseBool(isFinished); err == nil {
			p.Filter.IsFinished = &b
		}
	}

	notes, err := h.usecaseNote.FindAll(r.Context(), p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error("get notes: find all: db err: %v", err)
		return
	}

	json.NewEncoder(w).Encode(notes)
}

func(h *Note) GetNote(w http.ResponseWriter, r *http.Request) {
	noteId, err := getIdFromRequest(r, noteIdParam)
	if err != nil {
		writeErrorMessage(w, http.StatusBadRequest, wrongPathParams)
		return
	}

	u, ok := middleware.UserFromContext(r, h.log, "get note")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	zone := r.URL.Query().Get(timezoneQueryParam)

	timeZone, ok := h.checkZoneRule(r.Context(), zone, u.Id)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	note, err := h.usecaseNote.FindNote(r.Context(), noteId, u.Id, timeZone)
	if _, ok := err.(*usecase.ElemNotFound); ok {
		writeErrorMessage(w, http.StatusNotFound, noNoteFound)
		h.log.Warn(fmt.Sprintf("get note: note not found: user(id: %v) note(id: %v)", u.Id, noteId))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(fmt.Sprintf("get note: user(id: %v) note(id: %v) err: %v", u.Id, noteId, err))
		return
	}

	json.NewEncoder(w).Encode(note)
}

func(h *Note) PartialUpdateNote(w http.ResponseWriter, r *http.Request) {
	var n dto.NoteUpdate
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeErrorMessage(w, http.StatusBadRequest, wrongBody)
		return
	}

	noteId, err := getIdFromRequest(r, noteIdParam)
	if err != nil {
		writeErrorMessage(w, http.StatusBadRequest, wrongPathParams)
		return
	}

	u, ok := middleware.UserFromContext(r, h.log, "patch note update")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	note := model.NewNote(noteId, u.Id, n.Title, n.Text, n.Date, n.IsFinished)

	err = h.usecaseNote.UpdateNote(r.Context(), note)
	if _, ok := err.(*usecase.ElemNotFound); ok {
		w.WriteHeader(http.StatusNotFound)
		h.log.Warn(fmt.Sprintf("update note: note not found: user(id: %v) note(id: %v)", u.Id, noteId))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(fmt.Sprintf("update note: user(id: %v) note(id: %v) err: %v", u.Id, noteId, err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func(h *Note) DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteId, err := getIdFromRequest(r, noteIdParam)
	if err != nil {
		writeErrorMessage(w, http.StatusBadRequest, wrongPathParams)
		return
	}

	usr, ok := middleware.UserFromContext(r, h.log, "delete note")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.usecaseNote.RemoveNote(r.Context(), noteId, usr.Id)
	if _, ok := err.(*usecase.ElemNotFound); ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(fmt.Sprintf("update note: user(id: %v) note(id: %v) err: %v", usr.Id, noteId, err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func(h *Note) checkZoneRule(ctx context.Context, zone string, uId model.Id) (model.TimeZone, bool) {
	timeZone, ok := usecase.ValidateZone(zone)
	if ok {
		return timeZone, ok
	}

	usr, err := h.usecaseUser.FindById(ctx, uId)
	if err != nil {
		h.log.Error(fmt.Sprintf("get notes: checkZone: getUser: db err %v", err))
		return "", false
	}

	return usr.TimeZone, true
}
