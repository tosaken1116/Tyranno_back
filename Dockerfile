FROM golang:1.21.1-alpine AS dev

RUN apk update &&  apk add git
RUN go install github.com/cosmtrek/air@latest
WORKDIR /opt/nnyd

CMD ["air", "-c", ".air.toml"]

FROM golang:1.21.1-alpine AS prod

RUN mkdir -p /opt/nnyd

COPY . /opt/nnyd
RUN cat .env.sample | sed -e 's/ENV=.*/ENV=production/g' | sed 's/PORT=.*//g' > .env

WORKDIR /opt/nnyd

RUN go build ./cmd/main.go

CMD [ "./main" ]
