package store

type Store interface {
	User() UserRepository
	Redirect() RedirectRepository
}
