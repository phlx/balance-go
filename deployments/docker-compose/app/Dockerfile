FROM golang:1.15-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go build -mod=vendor -o service cmd/service/main.go

EXPOSE 8080

CMD ["./service"]
