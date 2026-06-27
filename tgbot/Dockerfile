# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS build

WORKDIR /src/tgbot

COPY tgbot/go.mod tgbot/go.sum* ./
RUN go mod download

COPY tgbot/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/tgbot ./cmd/tgbot

FROM alpine:3.22

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=build /out/tgbot /app/tgbot

USER app:app

ENTRYPOINT ["/app/tgbot"]
