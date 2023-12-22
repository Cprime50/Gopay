FROM golang:alpine as builder

WORKDIR /go/src/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /run .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /go/src/app/run .

EXPOSE 8082

CMD ["./run"]