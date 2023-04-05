package teststore_test

import (
	"testing"

	"github.com/leandoerbore/redirects/internal/app/model/redirect"
	"github.com/leandoerbore/redirects/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestRedirectRepository_Create(t *testing.T) {
	s := teststore.New()
	rdir := redirect.TestRedirect(t)

	assert.NoError(t, s.Redirect().Create(rdir))
	assert.NotNil(t, rdir)
}

func TestRedirectRepository_Find(t *testing.T) {
	s := teststore.New()
	rdir := redirect.TestRedirect(t)
	s.Redirect().Create(rdir)
	rdir, err := s.Redirect().Find(rdir.ID)
	assert.NoError(t, err)
	assert.NotNil(t, rdir)
}
