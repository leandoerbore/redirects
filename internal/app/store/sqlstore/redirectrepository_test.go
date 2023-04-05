package sqlstore_test

import (
	"testing"

	"github.com/leandoerbore/redirects/internal/app/model/redirect"
	"github.com/leandoerbore/redirects/internal/app/store"
	"github.com/leandoerbore/redirects/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestRedirectRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("redirects")

	s := sqlstore.New(db)
	rdir := redirect.TestRedirect(t)

	assert.NoError(t, s.Redirect().Create(rdir))
}

func TestRedirectRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("redirects")

	s := sqlstore.New(db)

	id := 100
	_, err := s.Redirect().Find(id)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	rdir := redirect.TestRedirect(t)
	rdir.ID = id
	s.Redirect().Create(rdir)
	rdir, err = s.Redirect().Find(rdir.ID)
	assert.NoError(t, err)
	assert.NotNil(t, rdir)
}

func TestRedirectRepository_GetAll(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("redirects")

	s := sqlstore.New(db)

	rdir := redirect.TestRedirect(t)
	s.Redirect().Create(rdir)

	rdirs, err := s.Redirect().GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, rdirs)
}

func TestRedirectRepository_Update(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("redirects")

	s := sqlstore.New(db)

	rdir := redirect.TestRedirect(t)
	s.Redirect().Create(rdir)
	rdir.IsActive = false

	err := s.Redirect().Update(rdir)
	assert.NoError(t, err)

	findedRdir, err := s.Redirect().Find(rdir.ID)
	assert.NoError(t, err)
	assert.NotNil(t, findedRdir)
	findedRdir.ID = rdir.ID
	assert.Equal(t, rdir, findedRdir)
}
