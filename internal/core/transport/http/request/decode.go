package core_http_request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("invalid JSON: %w", core_errors.ErrInvalidArgument)
	}

	var err error

	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", formatValidationError(err), core_errors.ErrInvalidArgument)
	}

	return nil
}

func formatValidationError(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return err.Error()
	}

	msgs := make([]string, 0, len(ve))
	for _, fe := range ve {
		msgs = append(msgs, formatFieldError(fe))
	}
	return strings.Join(msgs, "; ")
}

func formatFieldError(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s: обязательное поле", field)
	case "min":
		return fmt.Sprintf("%s: минимум %s символов", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s: максимум %s символов", field, fe.Param())
	case "email":
		return fmt.Sprintf("%s: некорректный email", field)
	case "startswith":
		return fmt.Sprintf("%s: должно начинаться с '%s'", field, fe.Param())
	case "omitempty":
		return fmt.Sprintf("%s: некорректное значение", field)
	default:
		return fmt.Sprintf("%s: некорректное значение", field)
	}
}
