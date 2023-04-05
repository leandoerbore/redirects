package sqlstore

import (
	"database/sql"

	"github.com/leandoerbore/redirects/internal/app/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db                 *sql.DB
	userRepository     *UserRepository
	redirectRepository *RedirectRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Redirect() store.RedirectRepository {
	if s.redirectRepository != nil {
		return s.redirectRepository
	}

	s.redirectRepository = &RedirectRepository{
		store: s,
	}

	return s.redirectRepository
}
