package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	var validator = validator.New()
	if err := validator.Struct(dest); err != nil {
		return fmt.Errorf("validate json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
