FROM golang:1.16-alpine3.14

WORKDIR /app/backend

RUN apk update && apk upgrade

COPY . .

RUN go build

CMD ["./cats"]