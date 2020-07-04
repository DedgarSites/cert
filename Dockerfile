# /usr/local/bin/start.sh will start the service

FROM golang:latest 

# Pause indefinitely if asked to do so.
ARG PAUSE_ON_BUILD
RUN test "$PAUSE_ON_BUILD" = "true" && while sleep 10; do true; done || :

COPY scripts/ /usr/local/bin/

ENV GOBIN=/bin \
    GOPATH=/go

RUN go get github.com/DedgarSites/cert && \
    cd /go/src/github.com/DedgarSites/cert && \
    go install && \
    cd && \
    rm -rf /go

EXPOSE 8443

USER 1001

CMD ["/usr/local/bin/start.sh"]
