FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO (required for go-sqlite3)
RUN apk add --no-cache build-base sqlite-dev 

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Enable CGO during build for sqlite
ENV CGO_ENABLED=1 

RUN make build


FROM alpine:latest

WORKDIR /app

# Install runtime dependencies in final image
RUN apk add --no-cache sqlite tzdata

COPY --from=builder /app/bin/chatbot-backend /app/bin/chatbot-backend
COPY .env /app/.env 

EXPOSE 8080

CMD ["./bin/chatbot-backend"]