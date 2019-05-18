FROM golang
RUN go get github.com/mattjw79/aws-api
ENTRYPOINT /go/bin/aws-api
EXPOSE 8080