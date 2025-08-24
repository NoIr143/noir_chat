build:
	@go build -o bin/noir_chat src/main.go

run: build
	@./bin/noir_chat

database:
	@docker exec -it noirdb psql -U noiruser -d noir_chat_db