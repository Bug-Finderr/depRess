FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o depres .

FROM alpine:latest
RUN apk add --no-cache graphviz
WORKDIR /output
COPY --from=builder /app/depres /usr/local/bin/depres
ENTRYPOINT ["depres"]
