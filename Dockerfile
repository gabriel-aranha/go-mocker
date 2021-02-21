FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go test -cover -v ./...
RUN go get -d -v ./...
RUN go install -v ./...

# Set Redis Username and Password
#ENV REDIS_HOST=""
#ENV REDIS_PASS=""

EXPOSE 1323

CMD ["/go/bin/go-mocker"]