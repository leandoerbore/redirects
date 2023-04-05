package redirect

import "testing"

func TestRedirect(t *testing.T) *Redirect {
	return &Redirect{
		Source:      "http://localhost:8000/home",
		Destination: "http://localhost:8000/newhome",
		IsActive:    true,
	}
}
