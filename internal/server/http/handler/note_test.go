package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todoNote/internal/model"
	"todoNote/internal/server/http/dto"
	"todoNote/internal/server/http/handler/mocks"
	"todoNote/internal/server/http/middleware"
	"todoNote/internal/usecase"
)

//go:generate mockgen -package=mocks -destination=mocks/user.go todoNote/internal/usecase INoteUsecase,IUserUsecase
//go:generate mockgen -package=mocks -destination=mocks/log.go todoNote/internal/server/http/log Logger


func TestNote_CreateNote(t *testing.T) {
	t.Run("", func(t *testing.T) {
		date := time.Now()
		b := dto.NewNote{Text: "text", Title: "title", Date: date}

		js, _ := json.Marshal(b)
		req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(js))
		q := req.URL.Query()
		q.Add(timezoneQueryParam, model.UTC)
		req.URL.RawQuery = q.Encode()

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockINoteUsecase(ctr)
		mockCase.EXPECT().CreateNote(gomock.Any(), gomock.Any()).
			Return(model.Id(1), nil).
			Do(func(_ context.Context, note *model.Note) {
				assert.Equal(t, model.Id(2), note.UserId)
				assert.Equal(t, false, note.IsFinished)
				assert.Equal(t, date.Format(time.RFC3339), note.Date.Format(time.RFC3339))
				assert.Equal(t, "text", note.Text)
				assert.Equal(t, "title", note.Title)
		})

		h := Note{usecaseNote: mockCase}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/notes", h.CreateNote)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 2})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		var note dto.IdObject
		json.NewDecoder(rr.Body).Decode(&note)
		fmt.Println(rr.Body)
		assert.Equal(t, model.Id(1), note.Id)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestNote_GetNote(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/notes/1", nil)
		q := req.URL.Query()
		q.Add(timezoneQueryParam, model.UTC)
		req.URL.RawQuery = q.Encode()

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockINoteUsecase(ctr)
		mockCase.EXPECT().FindNote(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&model.Note{Id: 1}, nil).
			Do(func(ctx context.Context, noteId model.Id, userId model.Id, zone model.TimeZone) {
				assert.Equal(t, model.Id(1), noteId, model.UTCp1)
		})

		h := Note{usecaseNote: mockCase}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/notes/{noteId}", h.GetNote)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 2})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		var note model.Note
		json.NewDecoder(rr.Body).Decode(&note)
		fmt.Println(rr.Body)
		assert.Equal(t, model.Id(1), note.Id)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("test failure", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/notes/15", nil)

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockINoteUsecase(ctr)
		mockCase.EXPECT().FindNote(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, usecase.NewNoteNotFoundError(15, 1)).
			Do(func(ctx context.Context, noteId model.Id, userId model.Id, zone model.TimeZone) {
				assert.Equal(t, model.Id(15), noteId, model.UTCp1)
			})

		mockUserCase := mocks.NewMockIUserUsecase(ctr)
		mockUserCase.EXPECT().FindById(gomock.Any(), gomock.Any()).
			Return(&model.User{Id: 1, TimeZone: model.UTCp3}, nil).
			Do(func(ctx context.Context, uId model.Id ) {
				assert.Equal(t, model.Id(1), uId)
			})

		mockLog := mocks.NewMockLogger(ctr)
		mockLog.EXPECT().Warn(gomock.Any())

		h := NewNoteHandler(mockCase, mockUserCase, mockLog)

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/notes/{noteId}", h.GetNote)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 1})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		var note model.Note
		json.NewDecoder(rr.Body).Decode(&note)
		fmt.Println(rr.Body)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestNote_GetNotes(t *testing.T) {
	t.Run("test with timezone query param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "api/v1/notes", nil)
		q := req.URL.Query()
		q.Add(timezoneQueryParam, model.UTC)
		req.URL.RawQuery = q.Encode()

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockINoteUsecase(ctr)
		mockCase.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return([]model.Note{
			{Id: 1, UserId: 1, Title: "title", Text: "text"},
			{Id: 2, UserId: 1, Title: "title", Text: "text"},
			{Id: 3, UserId: 1, Title: "title", Text: "text"}},
			nil).
			Do(func(ctx context.Context, filter usecase.FindParams) {
				assert.Equal(t, model.Id(1), filter.Filter.UserId)
				assert.Equal(t, model.UTC, filter.Zone)
		})

		h := Note{usecaseNote: mockCase}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.GetNotes)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 1})
		req = req.WithContext(ctx)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("test without timezone query param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "api/v1/notes", nil)

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockNoteCase := mocks.NewMockINoteUsecase(ctr)

		mockNoteCase.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return([]model.Note{
			{Id: 1, UserId: 1, Title: "title", Text: "text", Date: time.Now().UTC()},
			{Id: 2, UserId: 1, Title: "title", Text: "text", Date: time.Now().UTC()},
			{Id: 3, UserId: 1, Title: "title", Text: "text", Date: time.Now().UTC()}},
			nil).
			Do(func(ctx context.Context, filter usecase.FindParams) {
				assert.Equal(t, model.Id(1), filter.Filter.UserId)
				assert.Equal(t, model.UTCp3, filter.Zone)
			})

		mockUserCase := mocks.NewMockIUserUsecase(ctr)
		mockUserCase.EXPECT().FindById(gomock.Any(), gomock.Any()).
			Return(&model.User{Id: 1, TimeZone: model.UTCp3}, nil).
			Do(func(ctx context.Context, uId model.Id ) {
				assert.Equal(t, model.Id(1), uId)
		})

		h := Note{usecaseNote: mockNoteCase, usecaseUser: mockUserCase}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.GetNotes)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 1})
		req = req.WithContext(ctx)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var notes dto.Notes
		json.NewDecoder(rr.Body).Decode(&notes)
	})

	t.Run("test all query param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "api/v1/notes", nil)
		q := req.URL.Query()
		q.Add(limitQueryParam, "10")
		q.Add(offsetQueryParam, "10")
		now := time.Now()
		//so that we can check it in test, no need in nanos in backend
		now = time.UnixMilli(int64(now.UnixMilli() % 1000) * 1000)
		q.Add(startFromQueryParam, now.Format(time.RFC3339))
		q.Add(timezoneQueryParam, model.UTC)
		q.Add(isFinishedQueryParam, "false")
		req.URL.RawQuery = q.Encode()

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockINoteUsecase(ctr)
		mockCase.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return([]model.Note{
			{Id: 1, UserId: 1, Title: "title", Text: "text"},
			{Id: 2, UserId: 1, Title: "title", Text: "text"},
			{Id: 3, UserId: 1, Title: "title", Text: "text"}},
			nil).
			Do(func(ctx context.Context, filter usecase.FindParams) {
				assert.Equal(t, model.Id(1), filter.Filter.UserId)
				assert.Equal(t, model.UTC, filter.Zone)
				assert.Equal(t, uint64(10), *filter.Filter.Page.Offset)
				assert.Equal(t, uint64(10), *filter.Filter.Page.Limit)
				assert.Equal(t, now, *filter.Filter.TakeFrom)
				assert.Equal(t, false, *filter.Filter.IsFinished)
		})

		h := Note{usecaseNote: mockCase}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.GetNotes)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 1})
		req = req.WithContext(ctx)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}


func TestNote_PartialUpdateNote(t *testing.T) {
	t.Run("note exists", func(t *testing.T) {
		date := time.Now()
		b := dto.NoteUpdate{Title: "new title", Text: "new text", Date: date, IsFinished: true}

		js, _ := json.Marshal(b)
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/notes/2", bytes.NewReader(js))

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockINoteUsecase(ctr)
		mockCase.EXPECT().UpdateNote(gomock.Any(), gomock.Any()).
			Return(nil).
			Do(func(_ context.Context, note *model.Note) {
			assert.Equal(t, model.Id(2), note.UserId)
			assert.Equal(t, true, note.IsFinished)
			assert.Equal(t, date.Format(time.RFC3339), note.Date.Format(time.RFC3339))
			assert.Equal(t, "new text", note.Text)
			assert.Equal(t, "new title", note.Title)

		})

		h := Note{usecaseNote: mockCase}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/notes/{noteId}", h.PartialUpdateNote)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 2})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("note does not exits", func(t *testing.T) {
		date := time.Now()
		b := dto.NoteUpdate{Title: "new title", Text: "new text", Date: date, IsFinished: true}

		js, _ := json.Marshal(b)
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/notes/2", bytes.NewReader(js))

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockINoteUsecase(ctr)
		mockCase.EXPECT().UpdateNote(gomock.Any(), gomock.Any()).
			Return(usecase.NewNoteNotFoundError(2, 2)).
			Do(func(_ context.Context, note *model.Note) {
				assert.Equal(t, model.Id(2), note.UserId)
				assert.Equal(t, true, note.IsFinished)
				assert.Equal(t, date.Format(time.RFC3339), note.Date.Format(time.RFC3339))
				assert.Equal(t, "new text", note.Text)
				assert.Equal(t, "new title", note.Title)

			})

		mockLog := mocks.NewMockLogger(ctr)
		mockLog.EXPECT().Warn(gomock.Any())

		h := Note{usecaseNote: mockCase, log: mockLog}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/notes/{noteId}", h.PartialUpdateNote)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 2})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

}


func TestNote_DeleteNote(t *testing.T) {
	t.Run("note exists", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/notes/2", nil)

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockNoteCase := mocks.NewMockINoteUsecase(ctr)

		mockNoteCase.EXPECT().RemoveNote(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).
			Do(func(_ context.Context, nId, uId model.Id) {
				assert.Equal(t, model.Id(1), uId)
				assert.Equal(t, model.Id(2), nId)
		})

		h := Note{usecaseNote: mockNoteCase}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/notes/{noteId}",h.DeleteNote)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 1})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("note doesn't exist", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/notes/2", nil)

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockNoteCase := mocks.NewMockINoteUsecase(ctr)

		mockNoteCase.EXPECT().RemoveNote(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(usecase.NewNoteNotFoundError(2,1)).
			Do(func(_ context.Context, nId, uId model.Id) {
				assert.Equal(t, model.Id(1), uId)
				assert.Equal(t, model.Id(2), nId)
			})

		h := Note{usecaseNote: mockNoteCase}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/notes/{noteId}",h.DeleteNote)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 1})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})
}
