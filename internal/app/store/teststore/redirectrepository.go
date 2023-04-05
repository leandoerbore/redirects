package teststore

import (
	"github.com/leandoerbore/redirects/internal/app/model/redirect"
	"github.com/leandoerbore/redirects/internal/app/store"
)

type RedirectRepository struct {
	store     *Store
	redirects map[int]*redirect.Redirect
}

func (r *RedirectRepository) Create(rdir *redirect.Redirect) error {
	if err := rdir.Validate(); err != nil {
		return err
	}

	rdir.ID = len(r.redirects) + 1
	r.redirects[rdir.ID] = rdir

	return nil
}

func (r *RedirectRepository) Find(id int) (*redirect.Redirect, error) {
	rdir, ok := r.redirects[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return rdir, nil
}

func (r *RedirectRepository) GetAll() (*[]redirect.Redirect, error) {
	rdirs := make([]redirect.Redirect, 0)
	for _, rdir := range r.redirects {
		rdirs = append(rdirs, *rdir)
	}

	return &rdirs, nil
}

func (r *RedirectRepository) Update(rdir *redirect.Redirect) error {
	id := rdir.ID
	r.redirects[id] = rdir
	return nil
}

func (r *RedirectRepository) Remove(id int) error {
	if _, ok := r.redirects[id]; ok {
		delete(r.redirects, id)
	}
	return nil
}
