package in_memory

import (
	"context"
	"sort"
	"sync"
	"time"
	"todoNote/internal/model"
	"todoNote/internal/repo"
)

var _ repo.IRepoNote = &RepoNote{}

type RepoNote struct {
	sync.RWMutex
	storage map[model.Id]model.Note
	counter int64 //can make int64 and incr
}

func NewRepoNote() repo.IRepoNote {
	r := RepoNote{
		storage: make(map[model.Id]model.Note),
		counter: 1,
	}

	r.Insert(context.Background(), model.NewNote(1, 1, "title", "text", time.Now(), false))
	r.Insert(context.Background(), model.NewNote(2, 1, "title", "text", time.Now(), false))
	r.Insert(context.Background(), model.NewNote(3, 1, "title", "text", time.Now(), false))

	return &r
}

func(r *RepoNote) Insert(_ context.Context, n *model.Note) (model.Id, error) {
	n.Id = r.counter

	r.Lock()
	r.storage[n.Id] = *n
	r.counter++
	r.Unlock()

	return n.Id, nil
}

func(r *RepoNote) GetById(_ context.Context, id model.Id) (model.Note, error) {
	r.RLock()
	elem, ok := r.storage[id]
	r.RUnlock()
	if !ok {
		return model.Note{}, NewNoSuchElementError(id)
	}

	return elem, nil
}


type FindParams struct {
	Limit int
	Offset int
	LaterThan time.Time
	IsFinished bool
}

func(r *RepoNote) GetAllOffset(_ context.Context, filter repo.NoteFilter) ([]model.Note, error){
	filtered := make([]model.Note, 0)
	keys := make([]int, 0)

	r.RLock()
	for k, v := range r.storage {
		if v.UserId == filter.UserId {
			keys = append(keys, int(k))
		}
	}
	r.RUnlock()

	sort.Ints(keys)
	for i, k := range keys {
		if filter.Page.Offset != nil && uint64(i) < *filter.Page.Offset {
			continue
		}

		r.RLock()
		elem := r.storage[model.Id(k)]
		r.RUnlock()

		if filter.TakeFrom != nil && elem.Date.Before(*filter.TakeFrom) {
			continue
		}

		if filter.IsFinished != nil && elem.IsFinished != *filter.IsFinished {
			continue
		}

		filtered = append(filtered, elem)

		if filter.Page.Limit != nil && uint64(len(filtered)) >= *filter.Page.Limit {
			return filtered, nil
		}
	}


	return filtered, nil
}

func(r *RepoNote) Update(_ context.Context, n *model.Note) error {
	r.RLock()
	_, ok := r.storage[n.Id]
	r.RUnlock()
	//process
	if !ok  {
		return NewNoSuchElementError(n.Id)
	}

	r.Lock()
	r.storage[n.Id] = *n
	r.Unlock()
	return nil
}

func(r *RepoNote) Delete(_ context.Context, id model.Id) error {
	r.Lock()
	_, ok := r.storage[id]
	if ok {
		delete(r.storage, id)
	}
	r.Unlock()
	return nil
}
