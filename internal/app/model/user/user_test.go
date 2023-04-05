package user_test

import (
	"testing"

	"github.com/leandoerbore/redirects/internal/app/model/user"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *user.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *user.User {
				return user.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *user.User {
				u := user.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "email invalid",
			u: func() *user.User {
				u := user.TestUser(t)
				u.Email = "invalid"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *user.User {
				u := user.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *user.User {
				u := user.TestUser(t)
				u.Password = "123"
				return u
			},
			isValid: false,
		},
		{
			name: "with encrypt password",
			u: func() *user.User {
				u := user.TestUser(t)
				u.Password = ""
				u.EncryptedPassword = "encryptedpassword"
				return u
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	u := user.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}
