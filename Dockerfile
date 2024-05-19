FROM golang:1.22.2 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/main ./cmd/

ENV ENV="dev"

ENV GRPC_PORT=44045
ENV GRPC_TIMEOUT=1h

ENV MONGODB_URI="mongodb://admin:admin@localhost:27017"
ENV MONGODB_PING_TIMEOUT="10s"
ENV MONGODB_DATABASE_NAME="ucms-posts-dev"

ENV RABBITMQ_USER="admin"
ENV RABBITMQ_PASSWORD="admin"
ENV RABBITMQ_HOST="localhost"
ENV RABBITMQ_PORT="5672"


EXPOSE 44046

ENTRYPOINT ["./build/main"]