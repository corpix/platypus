FROM fedora:latest

RUN mkdir                       /etc/market-fetcher-http
ADD ./build/market-fetcher-http /usr/bin/market-fetcher-http

CMD [                                      \
    "/usr/bin/market-fetcher-http",        \
    "--config",                            \
    "/etc/market-fetcher-http/config.json" \
]
