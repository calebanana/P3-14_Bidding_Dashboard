# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY *.json ./

RUN go build -o /BiddingAPI

EXPOSE 8221

CMD [ "/BiddingAPI" ]