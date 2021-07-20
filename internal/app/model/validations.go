package model

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func notEqual(param bool) validation.RuleFunc {
	return func(value interface{}) error {
		v, _ := value.(bool)
		if v == param {
			return errors.New("unexpected bool")
		}

		return nil
	}
}
