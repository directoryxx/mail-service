# base go image
FROM golang:1.19-alpine as builder

RUN apk add build-base librdkafka-dev pkgconf

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=1 go build -tags musl -o ./build/mailApp ./internal/app

RUN chmod +x /app/build/mailApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

RUN mkdir /internal

COPY --from=builder /app/build/mailApp /app

COPY --from=builder /app/internal /internal

# COPY . /app

COPY ./.env /.env

CMD [ "/app/mailApp" ]