FROM golang:1.24.0-alpine3.20 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache make

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./src/cmd/main.go

FROM alpine:3.20

WORKDIR /app/

RUN apk add --no-cache curl 

COPY --from=builder /build/main .
COPY .env .
COPY ./src/config/configs .

CMD ["/app/main", "--config=default.yml"]