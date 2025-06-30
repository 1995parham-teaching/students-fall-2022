package request

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type StudentCreate struct {
	Name string `json:"name"`
}

func (r StudentCreate) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Length(1, 0), validation.Required),
	)
	if err != nil {
		return fmt.Errorf("student creation request validation failed %w", err)
	}

	err = validation.Validate(strings.Fields(r.Name),
		validation.Each(is.UTFLetter),
	)
	if err != nil {
		return fmt.Errorf("student creation request validation failed %w", err)
	}

	return nil
}
