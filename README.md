market-fetcher-http
---------

[![Build Status](https://travis-ci.org/cryptounicorns/market-fetcher-http.svg?branch=master)](https://travis-ci.org/cryptounicorns/market-fetcher-http)

HTTP interface for data providen by [market-fetcher](https://github.com/cryptounicorns/market-fetcher).

## Development

All development process accompanied by containers. Docker containers used for development, Rkt containers used for production.

> I am a big fan of Rkt, but it could be comfortable for other developers to use Docker for development and testing.

## Requirements

- [docker](https://github.com/moby/moby)
- [docker-compose](https://github.com/docker/compose)
- [jq](https://github.com/stedolan/jq)
- [rkt](https://github.com/coreos/rkt)
- [acbuild](https://github.com/containers/build)

### Preparations

#### For NSQ

> We use NSQ message queue by default but service also support kafka
> which is not used here because NSQ is quite enough.

##### Docker

``` console
$ sudo docker-compose up nsqd nsqlookupd nsqadmin

# Run a consumer for ticker topic
$ sudo docker run -it --net=host nsqio/nsq \
    /nsq_tail                              \
    --nsqd-tcp-address=127.0.0.1:4150      \
    --topic=ticker
```

##### Rkt

> All commands should be run in separate terminal windows.

``` console
$ sudo rkt run --interactive corpix.github.io/nsq:1.0.0 \
    --net=host                                          \
    -- nsqd                                             \
        --broadcast-address=127.0.0.1                   \
        --lookupd-tcp-address=127.0.0.1:4160            \
        --tcp-address=127.0.0.1:4150

$ sudo rkt run --interactive corpix.github.io/nsq:1.0.0 \
    --net=host                                          \
    -- nsqlookupd                                       \
        --tcp-address=127.0.0.1:4160

# Optionally you could run nsqadmin which will provide
# you a WEBUI for nsq topics etc.
$ sudo rkt run --interactive corpix.github.io/nsq:1.0.0 \
    --net=host                                          \
    -- nsqadmin                                         \
        --lookupd-http-address=127.0.0.1:4161           \
        --http-address=127.0.0.1:4171

# Run a consumer for ticker topic
$ sudo rkt run --interactive corpix.github.io/nsq:1.0.0 \
    --net=host                                          \
    -- nsq_tail                                         \
        --nsqd-tcp-address=127.0.0.1:4150               \
        --topic=ticker
```

#### Running market-fetcher

We need a data producer. It called a [market-fetcher](https://github.com/cryptounicorns/market-fetcher).

This is a service which fetches market data and puts it into the message queue our HTTP API will
be reading from.

Build a binary release:

``` console
$ git clone https://github.com/cryptounicorns/market-fetcher
$ cd market-fetcher
$ GOOS=linux make
```

##### Docker

> From the root of `market-fetcher` repository.

``` console
$ sudo docker-compose build market-fetcher
```

Now you should have a `cryptounicorns/market-fetcher` container. Start it:

> From the root of `market-fetcher-http` repository.

``` console
$ sudo docker-compose up market-fetcher
```

##### Rkt

There is no rkt container for this service at this time.

##### No isolation

> From the root of `market-fetcher` repository.

``` console
$ go run ./market-fetcher/market-fetcher.go --debug
```

### Running fetcher HTTP API frontend

Build a binary release:

``` console
$ GOOS=linux make
# This will put a binary into ./build/market-fetcher-http
```

#### Docker

``` console
$ docker-compose up market-fetcher-http
```

#### Rkt

There is no rkt container for this service at this time.

#### No isolation

``` console
$ go run ./market-fetcher-http/market-fetcher-http.go --debug
```

## Testing

> wsd is https://github.com/alexanderGugel/wsd

``` console
$ wsd -url ws://127.0.0.1:8080/api/v1/tickers/stream
< {"High":45.000001,"Low":41.35,"Avg":43.175001,"Vol":13157.821,"VolCur":302.6894,"Last":44.7,"Buy":44.749501,"Sell":44.978,"Timestamp":1501939951,"CurrencyPair":"LTC-USD","Market":"yobit"}
< {"High":2720,"Low":2504,"Avg":2612,"Vol":1608422.4,"VolCur":616.98825,"Last":2720,"Buy":2670,"Sell":2720,"Timestamp":1501939978,"CurrencyPair":"LTC-RUB","Market":"yobit"}
...
```
