FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/demo-app ./cmd/app

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /bin/demo-app /usr/local/bin/demo-app

EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/demo-app"]
