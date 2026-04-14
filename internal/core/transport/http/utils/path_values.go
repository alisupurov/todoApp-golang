package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf("no key=%s path value found: %w", key, core_errors.ErrInvalidArgument)
	}

	value, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf("invalid value for key=%s: %v: %w", key, err, core_errors.ErrInvalidArgument)
	}
	return value, nil
}