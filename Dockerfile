FROM golang:1.9.2 as builder

WORKDIR /go/src/github.com/cryptounicorns/platypus
COPY    . .
RUN     make

FROM fedora:latest

RUN  mkdir          /etc/platypus
COPY                /go/src/github.com/cryptounicorns/platypus/config.json     /etc/platypus/config.json
COPY --from=builder /go/src/github.com/cryptounicorns/platypus/build/platypus  /usr/bin/platypus

CMD [                           \
    "/usr/bin/platypus",        \
    "--config",                 \
    "/etc/platypus/config.json" \
]
