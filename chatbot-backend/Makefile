build:
	@go build -o bin/chatbot-backend main.go

testfast:
	@go test -v ./... -failfast

test:
	@go test -v ./...

run: build
	@./bin/chatbot-backend