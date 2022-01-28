package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
	"todoNote/internal/model"
	"todoNote/internal/server/http/dto"
	"todoNote/internal/server/http/handler/mocks"
	"todoNote/internal/server/http/middleware"
)

//go:generate mockgen -package=mocks -destination=mocks/user.go todoNote/internal/usecase IUserUsecase
//go:generate mockgen -package=mocks -destination=mocks/log.go todoNote/internal/server/http/log Logger
//go:generate mockgen -package=mocks -destination=mocks/auth.go todoNote/internal/server/http/auth IAuth


func TestAuth_Login(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		b := dto.UserLogin{UserName: "user", Password: "123"}
		js, _ := json.Marshal(b)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(js))

		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockCase := mocks.NewMockIUserUsecase(ctr)
		hash, _ := bcrypt.GenerateFromPassword([]byte("123"), 14)
		mockCase.EXPECT().FindByName(gomock.Any(), gomock.Any()).
			Return(&model.User{Id: 1, Name: "user", PasswordHash: hash}, nil).
			Do(func(_ context.Context, name string){
				assert.Equal(t, "user", name)
		})

		mockAuth := mocks.NewMockIAuth(ctr)
		mockAuth.EXPECT().CreateToken(gomock.Any()).
			Return(dto.JwtToken("token"), nil).
			Do(func(u model.UserInReq) {
				assert.Equal(t, model.Id(1), u.Id)
		})

		h := Auth{usecaseUser: mockCase, auth: mockAuth}

		rr := httptest.NewRecorder()
		ch := chi.NewRouter()
		ch.HandleFunc("/api/v1/login", h.Login)
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserAuthorized, model.UserInReq{Id: 1})
		req = req.WithContext(ctx)
		ch.ServeHTTP(rr, req)

		var tk dto.Token
		json.NewDecoder(rr.Body).Decode(&tk)
		fmt.Println(rr.Body)
		assert.Equal(t, "token", tk.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}


