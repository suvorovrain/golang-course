FROM golang:1.25

COPY . /tests
WORKDIR /tests

RUN go mod download

ENTRYPOINT ["go", "test", "-v", "./..."]