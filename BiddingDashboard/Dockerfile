# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY ./web ./web
COPY *.json ./

RUN go build -o /BiddingDashboard

EXPOSE 8220

CMD [ "/BiddingDashboard" ]