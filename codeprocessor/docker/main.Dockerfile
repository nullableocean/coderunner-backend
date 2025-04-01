FROM golang:1.24.0-alpine3.20 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./src/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /build/app .
COPY .env .
COPY docker/coderunner .

ENTRYPOINT ["/app/app", "--runnerdoc=/app/coderunner"]