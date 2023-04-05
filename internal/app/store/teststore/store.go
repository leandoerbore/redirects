package teststore

import (
	"github.com/leandoerbore/redirects/internal/app/model/redirect"
	"github.com/leandoerbore/redirects/internal/app/model/user"
	"github.com/leandoerbore/redirects/internal/app/store"
	_ "github.com/lib/pq"
)

type Store struct {
	userRepository     *UserRepository
	redirectRepository *RedirectRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*user.User),
	}

	return s.userRepository
}

func (s *Store) Redirect() store.RedirectRepository {
	if s.redirectRepository != nil {
		return s.redirectRepository
	}

	s.redirectRepository = &RedirectRepository{
		store:     s,
		redirects: make(map[int]*redirect.Redirect),
	}

	return s.redirectRepository
}
