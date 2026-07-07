package web_fs_repository

import (
	"fmt"
	"os"

	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
)

// GetFile читает файл по пути filePath и возвращает его содержимое как domain.File.
// Преобразует os.ErrNotExist в core_errors.ErrNotFound для единообразной обработки
// на уровне транспортного слоя (→ HTTP 404).
func (r *WebRepository) GetFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"file: %s: %w",
				filePath,
				core_errors.ErrNotFound,
			)
		}

		return nil, fmt.Errorf("get file: %s: %w", filePath, err)
	}

	return file, nil
}
