FROM golang:1.14

WORKDIR $GOPATH/src/github.com/TheLindaProjectInc/janus
COPY . $GOPATH/src/github.com/TheLindaProjectInc/janus
RUN go get -d ./...

CMD [ "go", "test", "-v", "./..."]