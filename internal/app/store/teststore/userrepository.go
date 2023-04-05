package teststore

import (
	"github.com/leandoerbore/redirects/internal/app/model/user"
	"github.com/leandoerbore/redirects/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*user.User
}

func (r *UserRepository) Create(u *user.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = len(r.users) + 1
	r.users[u.ID] = u

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (r *UserRepository) Find(id int) (*user.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}
