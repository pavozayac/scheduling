FROM golang:1.22.5

ENV GOOSE_DBSTRING postgresql://user:12345678@localhost:5432/constraint_db
ENV GOOSE_MIGRATION_DIR ./persistence/migrations
ENV GOOSE_DRIVER postgres

WORKDIR /app

COPY ./ ./

RUN ls -lR

RUN go install github.com/pressly/goose/v3/cmd/goose
# RUN goose up

RUN CGO_ENABLED=0 GOOS=linux make build

EXPOSE 8080

CMD [/app/server]