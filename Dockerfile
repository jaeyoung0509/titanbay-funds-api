FROM golang:1.26.2-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/titanbay-api ./cmd/api

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/titanbay-api /usr/local/bin/titanbay-api
COPY migrations ./migrations
COPY seed ./seed
COPY docs ./docs

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/titanbay-api"]
