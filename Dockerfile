FROM golang:1.20-alpine

WORKDIR /opt

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

COPY cmdhandlers/ ./cmdhandlers/
COPY executors/ ./executors/
COPY database/ ./database/

RUN go build -o ./tggoweatherbot

EXPOSE 8080

CMD [ "/opt/tggoweatherbot" ]