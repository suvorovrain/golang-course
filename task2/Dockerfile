FROM golang:1.26-alpine

ARG SERVICE_PATH
WORKDIR /app
RUN apk add --no-cache ca-certificates git
ENV GOPROXY=direct
ENV GOSUMDB=off
COPY . .
RUN go build -o /app/app_bin ./${SERVICE_PATH}/cmd/main.go
ENTRYPOINT ["/app/app_bin"]
