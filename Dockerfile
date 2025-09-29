FROM golang:1.25-alpine AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main ./cmd/main.go

FROM alpine:3.13

WORKDIR /

COPY --from=builder /main /main

ENTRYPOINT ["/main"]