FROM golang:1.14

RUN mkdir /go/src/go-api-project-layout
WORKDIR go/src/go-api-project-layout
COPY . .

RUN go build

CMD ["./go-api-project-layout"]