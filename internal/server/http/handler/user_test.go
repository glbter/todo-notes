package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"todoNote/internal/server/http/middleware"

	"net/http"
	"net/http/httptest"
	"testing"
	"todoNote/internal/model"
	"todoNote/internal/server/http/dto"
	"todoNote/internal/server/http/handler/mocks"
)

//go:generate mockgen -package=mocks -destination=mocks/user.go todoNote/internal/usecase INoteUsecase,IUserUsecase
//go:generate mockgen -package=mocks -destination=mocks/log.go todoNote/internal/server/http/log Logger

func TestUser_CreateUser(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		b := dto.UserRegistration{UserName: "user", Password: "123", TimeZone: model.UTCp3}

		js, _ := json.Marshal(b)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(js))

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockIUserUsecase(ctr)
		mockCase.EXPECT().Create(gomock.Any(), gomock.Any()).
			Return(model.Id(1), nil).
			Do(func(_ context.Context, u *model.UserNew) {
				assert.Equal(t, "user", u.Name)
				assert.Equal(t, "123", u.Password)
				assert.Equal(t, model.UTCp3, u.TimeZone)
		})

		h := User{userCase: mockCase}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/users", h.CreateUser)
		ch.ServeHTTP(rr, req)

		var note dto.IdObject
		json.NewDecoder(rr.Body).Decode(&note)
		fmt.Println(rr.Body)
		assert.Equal(t, model.Id(1), note.Id)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("test bad request", func(t *testing.T) {
		b := dto.UserRegistration{UserName: "user", Password: "123", TimeZone: "utc+3"}

		js, _ := json.Marshal(b)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(js))

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockIUserUsecase(ctr)


		h := User{userCase: mockCase}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/users", h.CreateUser)
		ch.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestUser_PartialUpdateUser(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		b := dto.UserUpdate{TimeZone: model.UTCp3}
		js, _ := json.Marshal(b)
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/users", bytes.NewReader(js))

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockIUserUsecase(ctr)
		mockCase.EXPECT().Update(gomock.Any(), gomock.Any()).
			Return(nil).
			Do(func(_ context.Context, user model.UserUpdate) {
			assert.Equal(t, model.Id(2), user.Id)
			assert.Equal(t, model.UTCp3, user.TimeZone)
		})

		h := User{userCase: mockCase}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/users", h.PartialUpdateUser)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 2})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("test bad request", func(t *testing.T) {
		b := dto.UserUpdate{TimeZone: "utc+3"}
		js, _ := json.Marshal(b)
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/users", bytes.NewReader(js))

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockIUserUsecase(ctr)

		h := User{userCase: mockCase}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/users", h.PartialUpdateUser)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 2})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestUser_DeleteUser(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users", nil)

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockIUserUsecase(ctr)
		mockCase.EXPECT().Remove(gomock.Any(), gomock.Any()).
			Return(nil).
			Do(func(_ context.Context, id model.Id) {
				assert.Equal(t, model.Id(2), id)
		})

		h := User{userCase: mockCase}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/users", h.DeleteUser)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 2})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})
}


