FROM golang:latest

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o selltechllc .

EXPOSE 8080

CMD ["./selltechllc"]



