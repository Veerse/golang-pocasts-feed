FROM golang:1.14

RUN mkdir /go/src/podcast-feed-api
WORKDIR go/src/podcast-feed-api
COPY . .

RUN go build

CMD ["./podcast-feed-api"]