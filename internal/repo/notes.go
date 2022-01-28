package repo

import (
	"context"
	"strconv"
	"time"
	"todoNote/internal/model"
)

type IRepoNote interface {
	Insert(ctx context.Context, n *model.Note) (model.Id, error)
	GetById(ctx context.Context, noteId model.Id) (model.Note, error)
	GetAllOffset(ctx context.Context, filter NoteFilter) ([]model.Note, error)
	Update(ctx context.Context, n *model.Note) error
	Delete(ctx context.Context, noteId model.Id) error
}

type NoteFilter struct {
	Page PageFilter
	UserId model.Id
	TakeFrom *time.Time
	IsFinished *bool
}

type PageFilter struct {
	Limit *uint64
	Offset *uint64
}

func GetUIntParamPointer(p string) *uint64 {
	if p != "" {
		off, err := strconv.Atoi(p)
		if err != nil || off > 0 {
			return uIntAddress(uint64(off))
		}
	}

	return nil
}

func uIntAddress(x uint64) *uint64 { return &x}
