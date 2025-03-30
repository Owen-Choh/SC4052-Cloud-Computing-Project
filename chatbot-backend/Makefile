build:
	@go build -o bin/chatbot-backend main.go

testfast:
	@go test ./... -failfast

testsum:
	@go test ./...

testv:
	@go test -v ./...

run: build
	@./bin/chatbot-backend