FROM golang:1.24.2

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./internal ./internal

COPY config.yaml ./

EXPOSE 8080


CMD ["go", "run", "./cmd"]