FROM golang

RUN mkdir /app

ADD . /app

WORKDIR /app

ARG service

RUN go build -o main ./cmd/$service

EXPOSE 8080

CMD ["/app/main"]