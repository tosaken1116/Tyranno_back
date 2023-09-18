FROM golang:1.21.1-alpine

RUN apk update &&  apk add git
RUN go install github.com/cosmtrek/air@latest
WORKDIR /opt/nnyd

CMD ["air", "-c", ".air.toml"]
