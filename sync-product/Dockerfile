FROM golang:1.21.1 as builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/main main.go
# CMD [ "app" ]
FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /bin
COPY --from=builder /usr/local/bin/main /bin
RUN chmod +x /bin/main
EXPOSE 2566
CMD ["/main"]