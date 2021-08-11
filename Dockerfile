FROM golang:1.14-alpine

RUN echo $GOPATH
RUN apk add --no-cache make gcc musl-dev git
WORKDIR $GOPATH/src/github.com/TheLindaProjectInc/janus
COPY ./ $GOPATH/src/github.com/TheLindaProjectInc/janus
RUN go install github.com/TheLindaProjectInc/janus/cli/janus

ENV METRIX_RPC=http://metrix:testpasswd@localhost:33851
ENV METRIX_NETWORK=regtest

EXPOSE 23889

ENTRYPOINT [ "janus" ]