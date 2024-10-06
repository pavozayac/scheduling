FROM golang:1.22.5-alpine as builder

WORKDIR /app

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux make build

FROM scratch

COPY --from=builder /app/server ./server

EXPOSE 8080

CMD ["./server"]