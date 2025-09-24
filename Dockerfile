FROM golang:alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/main.go

FROM alpine:3.13

WORKDIR /app

COPY --from=builder /app/main /app/main

RUN chmod +x /app/main

ENTRYPOINT ["/app/main"]