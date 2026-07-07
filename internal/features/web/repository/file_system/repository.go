// Package web_fs_repository реализует репозиторий для чтения статических файлов
// из файловой системы.
package web_fs_repository

type WebRepository struct{}

func NewWebRepository() *WebRepository {
	return &WebRepository{}
}
