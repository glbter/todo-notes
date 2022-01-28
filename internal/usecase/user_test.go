package usecase

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"todoNote/internal/model"
	in_memory "todoNote/internal/repo/in-memory"
	"todoNote/internal/usecase/mocks"
)

//go:generate mockgen -package=mocks -destination=mocks/users.go todoNote/internal/repo IRepoUser

func TestUserUsecase_Create(t *testing.T) {
	tts := []struct{
		in model.UserNew
		want model.User
		outErr error
	} {
		{
			in: model.UserNew{Name: "new", Password: "123", TimeZone: model.UTCp3},
			want: model.User{Name: "new", PasswordHash: []byte("123"), TimeZone: model.UTCp3},
			outErr: nil,
		},
		{
			in: model.UserNew{Name: "new", Password: "123"},
			want: model.User{Name: "new", PasswordHash: []byte("123"), TimeZone: model.UTC},
			outErr: nil,
		},
	}

	for _, tt := range tts {
		t.Run("", func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			mockRepo := mocks.NewMockIRepoUser(ctr)
			mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).
				Return(model.Id(0), nil).
				Do(func(ctx context.Context, user *model.User) {
					assert.Equal(t, tt.want.Name, user.Name)
					assert.Equal(t, tt.want.TimeZone, user.TimeZone)
					assert.NotEqual(t, "", string(user.PasswordHash))
			})

			uc := NewUserUsecase(mockRepo)
			_, err := uc.Create(context.Background(), &tt.in)
			assert.Equal(t, tt.outErr, err)
		})
	}
}

func TestUserUsecase_FindById(t *testing.T) {
	t.Run("successful find", func(t *testing.T) {
		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockRepo := mocks.NewMockIRepoUser(ctr)
		mockRepo.EXPECT().GetById(gomock.Any(), model.Id(1)).
			Return(&model.User{Id: 1}, nil)

		uc := NewUserUsecase(mockRepo)
		usr, err := uc.FindById(context.Background(), 1)
		assert.Equal(t, model.Id(1), usr.Id)
		assert.Nil(t, err)
	})

	t.Run("unsuccessful find", func(t *testing.T) {
		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockRepo := mocks.NewMockIRepoUser(ctr)
		mockRepo.EXPECT().GetById(gomock.Any(), model.Id(1)).
			Return(nil, in_memory.NewNoSuchElementError(1))

		uc := NewUserUsecase(mockRepo)
		_, err := uc.FindById(context.Background(), 1)
		assert.NotNil(t, err)
	})
}

func TestUserUsecase_FindByName(t *testing.T) {
	t.Run("successful find", func(t *testing.T) {
		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockRepo := mocks.NewMockIRepoUser(ctr)
		mockRepo.EXPECT().GetByUserName(gomock.Any(), "user name").
			Return(&model.User{Id: 1}, nil)

		uc := NewUserUsecase(mockRepo)
		usr, err := uc.FindByName(context.Background(), "user name")
		assert.Equal(t, model.Id(1), usr.Id)
		assert.Nil(t, err)
	})

	t.Run("unsuccessful find", func(t *testing.T) {
		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockRepo := mocks.NewMockIRepoUser(ctr)
		mockRepo.EXPECT().GetByUserName(gomock.Any(), "user name").
			Return(nil, in_memory.NewNoSuchNameError("user name"))

		uc := NewUserUsecase(mockRepo)
		_, err := uc.FindByName(context.Background(), "user name")
		assert.NotNil(t, err)
	})
}

func TestUserUsecase_Update(t *testing.T) {
	tts := []struct{
		in model.UserUpdate
		getByIdUser model.User
		getByIdErr error
	} {
		{
			in: model.UserUpdate{Id: 1, TimeZone: model.UTC},
			getByIdUser: model.User{Id: 1},
			getByIdErr: nil,
		},
		{
			in: model.UserUpdate{Id: 1},
			getByIdUser: model.User{Id: 1, TimeZone: model.UTC},
			getByIdErr: nil,
		},
		{
			in: model.UserUpdate{Id: 2},
			getByIdUser: model.User{},
			getByIdErr: in_memory.NewNoSuchElementError(2),
		},
	}

	for _, tt := range tts {
		t.Run("", func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			mockRepo := mocks.NewMockIRepoUser(ctr)
			mockRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).
				Return(&tt.getByIdUser, tt.getByIdErr)
			if tt.getByIdErr == nil {
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Do(func(ctx context.Context, user *model.User) {
						assert.Equal(t, tt.getByIdUser.TimeZone, user.TimeZone)
					})
			}

			uc := NewUserUsecase(mockRepo)
			err := uc.Update(context.Background(), tt.in)
			if err != nil {
				assert.Equal(t, fmt.Sprintf("update user: %v", tt.getByIdErr.Error()), err.Error())
			}
		})
	}
}

func TestUserUsecase_Remove(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ctr := gomock.NewController(t)
		defer ctr.Finish()
		mockRepo := mocks.NewMockIRepoUser(ctr)
		mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).
			Return(nil).
			Do(func(ctx context.Context, id model.Id) {
				assert.Equal(t, model.Id(1), id)
		})

		uc := NewUserUsecase(mockRepo)
		err := uc.Remove(context.Background(), 1)
		assert.Nil(t, err)
	})
}

func TestUserUsecase_HashUser(t *testing.T) {
	uc := UserUsecase{}
	word := "u"
	h, _ := bcrypt.GenerateFromPassword([]byte(word), 14)

	tts := []struct{
		in model.UserNew
		out model.User
		err error
	} {
		{model.UserNew{Password: word}, model.User{PasswordHash: h}, nil},
	}

	for _, tt := range tts {
		got, err := uc.HashUser(&tt.in)
		//assert.Equal(t, string(tt.out.PasswordHash), string(got.PasswordHash))
		assert.Equal(t, tt.err, err)
		err = bcrypt.CompareHashAndPassword(got.PasswordHash, []byte(tt.in.Password))
		assert.Equal(t, tt.err, err)
	}
}

