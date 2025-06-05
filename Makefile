MIGRATIONS_PATH =./cmd/migrate/migrations
DB_URI=postgres://admin:adminpassword@localhost:5432/gophersocial?sslmode=disable

test:
	@go test -v ./...

migration-create:
	@goose -dir $(MIGRATIONS_PATH) -s create $(name) sql

migrate-up:
	@goose -dir $(MIGRATIONS_PATH) postgres "$(DB_URI)" up

migrate-down:
	@goose -dir $(MIGRATIONS_PATH) postgres "$(DB_URI)" down

seed: 
	@go run cmd/migrate/seed/main.go