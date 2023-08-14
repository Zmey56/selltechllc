FROM golang:1.20 AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=build /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]
