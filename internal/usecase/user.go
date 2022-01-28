package usecase

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"todoNote/internal/model"
	"todoNote/internal/repo"
)
type IUserUsecase interface {
	Create(ctx context.Context, usr *model.UserNew) (model.Id, error)
	FindById(ctx context.Context, uId model.Id) (*model.User, error)
	FindByName(ctx context.Context, name string) (*model.User, error)
	Update(ctx context.Context, usr model.UserUpdate) error
	Remove(ctx context.Context, uId model.Id) error
}
var _ IUserUsecase = &UserUsecase{}

type UserUsecase struct {
	userRepo repo.IRepoUser
}

func NewUserUsecase(r repo.IRepoUser) *UserUsecase {
	return &UserUsecase{
		userRepo: r,
	}
}

func (u *UserUsecase) Create(ctx context.Context, usr *model.UserNew) (model.Id, error) {
	user, err := u.HashUser(usr)
	if err != nil {
		return -1, fmt.Errorf("create user.go: %w", err)
	}

	if user.TimeZone == "" {
		user.TimeZone = model.UTC
	}

	id, err := u.userRepo.Insert(ctx, user);
	if err != nil {
		return 0, fmt.Errorf("create user.go: %w", err)
	}
	user.Id = id

	return id, nil
}

func(u *UserUsecase) FindById(ctx context.Context, uId model.Id) (*model.User, error) {
	usr, err := u.userRepo.GetById(ctx, uId)
	if err != nil {
		return nil, fmt.Errorf("findById: %w", err)
	}

	return usr, nil
}

func (u *UserUsecase) FindByName(ctx context.Context, name string) (*model.User, error) {
	usr, err := u.userRepo.GetByUserName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("findByName: %w", err)
	}

	return usr, nil
}

func (u *UserUsecase) Update(ctx context.Context, usr model.UserUpdate) error {
	user, err := u.userRepo.GetById(ctx, usr.Id)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if usr.TimeZone != "" {
		user.TimeZone = usr.TimeZone
	}

	return u.userRepo.Update(ctx, user)
}

func (u *UserUsecase) Remove(ctx context.Context, uId model.Id) error {
	if err := u.userRepo.Delete(ctx, uId); err != nil {
		return fmt.Errorf("user.go remove: %w", err)
	}

	return nil
}

func(u *UserUsecase) HashUser(usr *model.UserNew) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 14)
	if err != nil {
		return nil, fmt.Errorf("hash user.go: %w", err)
	}

	user := model.NewUser(0, usr.Name, hash, usr.TimeZone)
	return user, nil
}