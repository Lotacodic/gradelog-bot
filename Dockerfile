# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gradelog-bot .

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gradelog-bot .

CMD ["./gradelog-bot"]
