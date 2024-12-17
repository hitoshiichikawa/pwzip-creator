FROM golang:1.23.3-alpine3.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o passwordzipper

EXPOSE 8080

CMD ["./passwordzipper"]