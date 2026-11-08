.PHONY: help migrate-create migrate-up migrate-down compose-up compose-down compose-logs compose-ps compose-logs-service compose-rebuild-server

-include .env
export

service_name ?= server

# первая команда (выводится по умолчанию при make)
help:
	@echo "Миграции:"
	@echo "  make migrate-create migrate_name=имя   создать пару файлов миграции"
	@echo "  make migrate-up                        применить миграции (с хоста)"
	@echo "  make migrate-down                      откатить одну миграцию"
	@echo ""
	@echo "Docker Compose:"
	@echo "  make compose-up                        поднять всё в фоне"
	@echo "  make compose-down                      остановить и удалить контейнеры (тома остаются)"
	@echo "  make compose-logs                      логи всех сервисов"
	@echo "  make compose-logs-service service_name=postgres   логи одного сервиса"
	@echo "  make compose-ps                        статус сервисов"
	@echo "  make compose-rebuild-server            пересобрать и перезапустить сервер"


# группа команд для golang-migrate/migrate
migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(migrate_name)
	
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down 1

	
# группа команд для docker compose
compose-up:
	docker compose up -d

compose-down:
	docker compose down

compose-logs:
	docker compose logs -f -t

compose-logs-service:
	docker compose logs -f -t $(service_name)
	
compose-ps:
	docker compose ps -a

compose-rebuild-server:
	docker compose up -d --build server