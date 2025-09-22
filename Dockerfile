FROM golang:alpine3.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o app ./cmd/main.go
#from alipne скопировать бинарь в новый контейнер, чтобы контейнер весил меньше
CMD ["./app"]