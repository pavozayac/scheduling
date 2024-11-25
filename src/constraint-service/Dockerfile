FROM golang:1.22.5-alpine AS builder

RUN apk add make

WORKDIR /app

COPY ./ ./

RUN go install github.com/go-task/task/v3/cmd/task@latest

RUN CGO_ENABLED=0 GOOS=linux task build

FROM scratch

COPY --from=builder /tmp/bin/server ./server

EXPOSE 8080

CMD ["./server"]