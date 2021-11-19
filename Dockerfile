FROM golang:1.16.3-alpine3.13
RUN apk --no-cache add nano procps \
  && rm -rf /var/cache/apk/*
WORKDIR /root
COPY . .
RUN go build -o main main.go
ENTRYPOINT /bin/sh