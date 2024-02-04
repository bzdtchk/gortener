FROM golang:1.20.4-alpine AS builder

ENV CGO_ENABLED=1

RUN apk add --update gcc musl-dev sqlite

WORKDIR /go/src/app

COPY . .

RUN go mod download

FROM builder AS dev

RUN go get github.com/githubnemo/CompileDaemon && go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build='go build main.go' --command='./main'

FROM builder AS prod

ENV GIN_MODE=release

RUN sqlite3 data/data.db ".databases"

RUN go build main.go

ENTRYPOINT ./main