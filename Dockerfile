FROM fedora:latest

RUN mkdir            /etc/platypus
ADD ./build/platypus /usr/bin/platypus

CMD [                           \
    "/usr/bin/platypus",        \
    "--config",                 \
    "/etc/platypus/config.json" \
]
