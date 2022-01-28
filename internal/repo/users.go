package repo

import (
	"context"
	"todoNote/internal/model"
)

type IRepoUser interface {
	Insert(ctx context.Context, u *model.User) (model.Id, error)
	GetByUserName(ctx context.Context, name string) (*model.User, error)
	GetById(ctx context.Context, uId model.Id) (*model.User, error)
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, uId model.Id) error
}
