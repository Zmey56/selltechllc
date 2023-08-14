FROM golang:latest

ENV GO111MODULE=on
ENV GOPATH=/

COPY ./ ./

RUN go mod download

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]



