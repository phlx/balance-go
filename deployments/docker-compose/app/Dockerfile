FROM golang:1.15-alpine

ARG TZ='Europe/Moscow'

ENV DEFAULT_TZ ${TZ}

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh tzdata

RUN cp /usr/share/zoneinfo/${DEFAULT_TZ} /etc/localtime && echo ${DEFAULT_TZ} > /etc/timezone

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go build -mod=vendor -o service cmd/service/main.go

EXPOSE 8080

CMD ["./service"]
