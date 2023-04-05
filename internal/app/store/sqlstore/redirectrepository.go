package sqlstore

import (
	"database/sql"
	"errors"

	"github.com/leandoerbore/redirects/internal/app/model/redirect"
	"github.com/leandoerbore/redirects/internal/app/store"
)

type RedirectRepository struct {
	store *Store
}

var (
	errWithPatch = errors.New("something went wrong when patch")
)

func (r *RedirectRepository) Create(rdir *redirect.Redirect) error {
	if err := rdir.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO redirects (source, destination, is_active) VALUES ($1, $2, $3) RETURNING id",
		rdir.Source,
		rdir.Destination,
		rdir.IsActive,
	).Scan(&rdir.ID)
}

func (r *RedirectRepository) Find(id int) (*redirect.Redirect, error) {
	rdir := &redirect.Redirect{}

	if err := r.store.db.QueryRow(
		"SELECT id, source, destination, is_active FROM redirects WHERE id = $1",
		id,
	).Scan(
		&rdir.ID,
		&rdir.Source,
		&rdir.Destination,
		&rdir.IsActive,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return rdir, nil
}

func (r *RedirectRepository) GetAll() (*[]redirect.Redirect, error) {
	arrRdir := []redirect.Redirect{}

	rows, err := r.store.db.Query("SELECT id, source, destination, is_active  FROM redirects")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	for rows.Next() {
		rdir := redirect.Redirect{}

		if err := rows.Scan(&rdir.ID, &rdir.Source, &rdir.Destination, &rdir.IsActive); err != nil {
			return nil, err
		}
		arrRdir = append(arrRdir, rdir)
	}

	return &arrRdir, nil
}

func (r *RedirectRepository) Update(rdir *redirect.Redirect) error {
	rdirOld, err := r.store.redirectRepository.Find(rdir.ID)
	if err != nil {
		return err
	}
	if rdirOld == nil {
		return store.ErrRecordNotFound
	}

	r.store.db.QueryRow(
		`
		UPDATE redirects 
		SET source = $1, destination = $2, is_active = $3
		WHERE id = $4
		`,
		&rdir.Source,
		&rdir.Destination,
		&rdir.IsActive,
		&rdir.ID,
	)

	return nil
}

func (r *RedirectRepository) Remove(id int) error {
	rdir, err := r.store.redirectRepository.Find(id)
	if err != nil {
		return err
	}

	if rdir == nil {
		return nil
	}

	r.store.db.QueryRow(
		`
		DELETE FROM redirects
		WHERE id = $1
		`,
		id,
	)

	return nil
}
