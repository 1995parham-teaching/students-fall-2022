package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type StudentCreate struct {
	Name string `json:"name"`
}

func (r StudentCreate) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, is.UTFLetter, validation.Length(1, 0)),
	)
}
