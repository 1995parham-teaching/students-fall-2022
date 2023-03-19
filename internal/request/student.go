package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type StudentCreate struct {
	Name string `json:"name"`
}

func (r StudentCreate) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, is.UTFLetter, validation.Length(1, 0)),
	); err != nil {
		return fmt.Errorf("student creation request validation failed %w", err)
	}

	return nil
}
