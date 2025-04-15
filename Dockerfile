# syntax=docker/dockerfile:1

FROM golang:1.24.2

WORKDIR /server

COPY . .

RUN go build -o /gameserver cmd/server/main.go

EXPOSE 50051
