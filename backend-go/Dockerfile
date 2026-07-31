# backend-go/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/server ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/models ./models

EXPOSE 8080
CMD ["./server"]