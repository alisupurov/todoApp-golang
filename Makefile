include .env
export

export PROJECT_ROOT=${shell pwd}

.DEFAULT_GOAL := help

env-up: ## env: Запустить окружение проекта (PostgreSQL + Redis)
	@docker compose up -d todo-postgres todo-redis

env-down: ## env: Остановить окружение проекта
	@docker compose down todo-postgres todo-redis

env-cleanup: ## env: Очистить volume файлы окружения
	@read -p "Очистить все volume файлы? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todo-postgres port-forwarder && \
		sudo rm -rf ${PROJECT_ROOT}/.out/pgdata&& \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

migrate-create: ## PostgreSQL: Создать новую версию схемы данных
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр 'seq'. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm todo-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр 'action'. Пример: make migrate-action action=down 2"; \
		exit 1; \
	fi; \
	docker compose run --rm todo-postgres-migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todo-postgres:5432/${POSTGRES_DB}?sslmode=disable \
	"$(action)"

migrate-up: ## PostgreSQL: Накатить миграции
	@make migrate-action action=up

migrate-down: ## PostgreSQL: Откатить миграции
	@make migrate-action action=down

env-port-forward: ## env: Открыть порты сервисов окружения
	@docker compose up -d port-forwarder

env-port-close: ## env: Закрыть порты сервисов окружения
	@docker compose down port-forwarder

logs-cleanup: ## env: Очистить файлы логов
	@read -p "Очистить все log файлы? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todo-postgres port-forwarder && \
		sudo rm -rf ${PROJECT_ROOT}/.out/logs&& \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi

todoapp-run: ## Golang приложение: Запустить локально (для разработки)
	@export LOGGER_FOLDER=${PROJECT_ROOT}/.out/logs && \
	export POSTGRES_HOST=localhost && \
	export REDIS_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

todoapp-deploy: ## Golang приложение: Запустить в Docker Compose (для деплоя)
	@docker compose up -d --build todoapp

todoapp-undeploy: ## Golang приложение: Остановить Docker Compose сервис
	@docker compose down todoapp

swagger-gen: ## Swagger: Сгенерировать актуальную спецификацию
	@docker compose run --rm swagger \
		init \
		-g cmd/todoapp/main.go \
		-o docs \
		--parseInternal \
		--parseDependency
ps: ## env: Посмотреть запущенные Docker Compose сервисы
	@docker compose ps

help: ## Показать справку по командам
	@echo "=== Центр управления проектом ==="
	@echo ""
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
