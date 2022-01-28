package in_memory

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"todoNote/internal/model"
	"todoNote/internal/repo"
)

var _ repo.IRepoUser = &RepoUser{}

type RepoUser struct {
	sync.RWMutex
	storage map[model.Id]model.User
	counter int64
}

func NewRepoUser() repo.IRepoUser {
	r := RepoUser{
		storage: make(map[model.Id]model.User),
		counter: 1,
	}

	h, _ := bcrypt.GenerateFromPassword([]byte("user"), 14)
	r.Insert(context.Background(), model.NewUser(1, "user", h, model.UTCm4))

	return &r
}

func(r *RepoUser) Insert(_ context.Context, u *model.User) (model.Id, error) {
	u.Id = r.counter

	r.Lock()
	r.storage[u.Id] = *u
	r.counter++
	r.Unlock()

	return u.Id, nil
}

func(r *RepoUser) GetByUserName(_ context.Context, uName string) (*model.User, error) {
	r.RLock()
	defer r.RUnlock()

	for _, v := range r.storage {
		if v.Name == uName {
			return &v, nil
		}
	}

	return nil, NewNoSuchNameError(uName)
}

func(r *RepoUser) GetById(ctx context.Context, uId model.Id) (*model.User, error) {
	r.RLock()
	u, ok := r.storage[uId]
	r.RUnlock()
	if !ok {
		return &u, NewNoSuchElementError(uId)
	}

	return &u, nil
}


func(r *RepoUser) Update(_ context.Context, u *model.User) error {
		r.RLock()
		_, ok := r.storage[u.Id]
		r.RUnlock()

		if !ok {
			return NoSuchElementError{id: u.Id}
		}

		r.Lock()
		r.storage[u.Id] = *u
		r.Unlock()

		return nil
}

func(r *RepoUser) Delete(_ context.Context, userId model.Id) error {
	r.Lock()
	delete(r.storage, userId)
	r.Unlock()

	return nil
}
