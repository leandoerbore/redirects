package store

import (
	"github.com/leandoerbore/redirects/internal/app/model/redirect"
	"github.com/leandoerbore/redirects/internal/app/model/user"
)

type UserRepository interface {
	Create(*user.User) error
	Find(int) (*user.User, error)
	FindByEmail(string) (*user.User, error)
}

type RedirectRepository interface {
	Create(*redirect.Redirect) error
	Find(int) (*redirect.Redirect, error)
	GetAll() (*[]redirect.Redirect, error)
	Update(*redirect.Redirect) error
	Remove(int) error
}
