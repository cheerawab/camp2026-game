# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS build

WORKDIR /src/server

RUN apk add --no-cache ca-certificates curl

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./
ARG QUIZ_QUESTIONS_CSV_URL
RUN QUIZ_QUESTIONS_CSV_URL="${QUIZ_QUESTIONS_CSV_URL}" sh ./scripts/fetch-quiz-questions.sh
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.22

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=build /out/api /app/api
COPY --from=build /src/server/content /app/content

ENV CONTENT_DIR=/app/content

USER app:app

EXPOSE 8080

ENTRYPOINT ["/app/api"]
