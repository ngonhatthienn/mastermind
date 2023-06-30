FROM golang:1.18-alpine

WORKDIR /intern2023
COPY go.mod ./
COPY go.sum ./
RUN go get -d -v ./...

COPY . .

RUN go build -o main

EXPOSE 8090

CMD ["./main"]

