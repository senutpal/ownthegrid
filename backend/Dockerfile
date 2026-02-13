FROM golang:1.22-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /app/server ./server

EXPOSE 8080

USER app

ENTRYPOINT ["./server"]
