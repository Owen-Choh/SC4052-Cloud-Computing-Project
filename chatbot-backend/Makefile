build:
	@go build -o bin/chatbot-backend main.go

test:
	@go test ./...

testfast:
	@go test ./... -failfast

testv:
	@go test -v ./...

run: build
	@./bin/chatbot-backend