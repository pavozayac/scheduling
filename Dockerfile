FROM golang:1.24.0-alpine AS builder

WORKDIR /app

RUN apk add go-task-task

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux task build

FROM scratch

COPY --from=builder /tmp/bin/server ./server
COPY ./internal/infrastructure/db/migrations ./migrations

EXPOSE 8080

CMD ["./server"]