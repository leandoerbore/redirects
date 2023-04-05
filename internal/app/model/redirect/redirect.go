package redirect

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Redirect struct {
	ID          int    `json:"id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	IsActive    bool   `json:"is_active"`
}

func (rdir *Redirect) Validate() error {
	return validation.ValidateStruct(
		rdir,
		validation.Field(&rdir.Source, validation.Required, validation.Length(6, 100)),
		validation.Field(&rdir.Destination, validation.Required, validation.Length(6, 100)),
	)
}
