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
< {"high":2943.7997,"low":2745.7621,"avg":2844.7809,"vol":1584868.5,"volCur":564.85726,"last":2745.7621,"buy":2766.2631,"sell":2788.4115,"timestamp":1503190041,"currencyPair":"LTC-RUB","market":"yobit"}
< {"high":4350,"low":4116.798,"avg":0,"vol":34787.59018931,"volCur":1119.71769938,"last":4280.783,"buy":4280.8723,"sell":4285.3211,"timestamp":1503190091,"currencyPair":"BTC-USD","market":"cex"}
< {"high":266882.1,"low":252200.02,"avg":0,"vol":60.3862276,"volCur":2.35581065,"last":266882.1,"buy":256782.46,"sell":263975,"timestamp":1503190091,"currencyPair":"BTC-RUB","market":"cex"}
...
```
