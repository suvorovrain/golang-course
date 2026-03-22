FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/collector ./collector/cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /out/collector ./collector
EXPOSE 50051
ENTRYPOINT ["./collector"]
