FROM golang:1.19.2-alpine3.16

RUN apk update && \
    apk add --no-cache gcc musl-dev

WORKDIR /exbestfriend

COPY . .

RUN go mod download

WORKDIR /exbestfriend/cmd/instaspy

CMD ["go", "run", "main.go"]
