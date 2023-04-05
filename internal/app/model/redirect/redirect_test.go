package redirect_test

import (
	"testing"

	"github.com/leandoerbore/redirects/internal/app/model/redirect"
	"github.com/stretchr/testify/assert"
)

func TestRedirect_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		rdir    func() *redirect.Redirect
		isValid bool
	}{
		{
			name: "valid",
			rdir: func() *redirect.Redirect {
				return redirect.TestRedirect(t)
			},
			isValid: true,
		},
		{
			name: "empty from url",
			rdir: func() *redirect.Redirect {
				rdir := redirect.TestRedirect(t)
				rdir.Source = ""
				return rdir
			},
			isValid: false,
		},
		{
			name: "empty destination url",
			rdir: func() *redirect.Redirect {
				rdir := redirect.TestRedirect(t)
				rdir.Destination = ""
				return rdir
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.rdir().Validate())
			} else {
				assert.Error(t, tc.rdir().Validate())
			}
		})
	}
}
