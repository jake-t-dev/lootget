FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o bot ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates


WORKDIR /root/
COPY --from=builder /app/bot .

CMD ["./bot"]
