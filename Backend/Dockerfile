# Build stage
FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server


# Runtime stage
FROM alpine:latest

WORKDIR /app

# copy binary
COPY --from=builder /app/main .

# 🔥 COPY MIGRATIONS (IMPORTANT)
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]