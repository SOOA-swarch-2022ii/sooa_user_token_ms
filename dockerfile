# syntax=docker/dockerfile:1
FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . ./

RUN go build -o /sooa_user_token_ms

EXPOSE 6665

CMD [ "/sooa_user_token_ms" ]