FROM golang:1.19-alpine

WORKDIR /usr/src/app

EXPOSE 80

COPY . .

RUN go install github.com/cosmtrek/air@latest

CMD ["air"]
