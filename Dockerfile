FROM golang:1.19.2-alpine3.16

WORKDIR /exbestfriend

COPY . .

ENV CONFIG_PATH=/exbestfriend/cmd/instaspy/

ARG TELEGRAM_BOT
ARG CHAT_ID

ENV TELEGRAM_BOT=$TELEGRAM_BOT
ENV CHAT_ID=$CHAT_ID

RUN go mod download

WORKDIR /exbestfriend/cmd/instaspy

CMD ["go", "run", "main.go"]
