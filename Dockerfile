FROM golang:1.25-alpine AS builder

ENV GO111MODULE=on
ENV GOSUMDB=off

WORKDIR /
COPY go.mod go.sum ./
COPY vendor ./vendor

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /main ./cmd/main.go

FROM alpine:3.13
WORKDIR /
COPY --from=builder /main /main
ENTRYPOINT ["/main"]
