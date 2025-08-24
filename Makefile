build:
	@go build -o bin/noir_chat src/main.go

run: build
	@./bin/noir_chat

database:
	@docker exec -it noirdb psql -U noiruser -d noir_chat_db

migration:
	@migrate create -ext sql -dir src/database/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run src/database/migrate/main.go up

migrate-down:
	@go run src/database/migrate/main.go down