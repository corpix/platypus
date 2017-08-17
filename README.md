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

##### Docker

``` console
$ sudo docker-compose up nsqd nsqlookupd nsqadmin

# Run a consumer for ticker topic
$ sudo docker run -it --net=host nsqio/nsq /nsq_tail --nsqd-tcp-address=127.0.0.1:4150 --topic=ticker
```

##### Rkt

> All commands should be run in separate terminal windows.

``` console
$ sudo rkt run --interactive corpix.github.io/nsq:1.0.0 --net=host -- nsqd --broadcast-address=127.0.0.1 --lookupd-tcp-address=127.0.0.1:4160 --tcp-address=127.0.0.1:4150
$ sudo rkt run --interactive corpix.github.io/nsq:1.0.0 --net=host -- nsqlookupd --tcp-address=127.0.0.1:4160

# Optionally you could run nsqadmin which will provide you a WEBUI for nsq topics etc.
$ sudo rkt run --interactive corpix.github.io/nsq:1.0.0 --net=host -- nsqadmin --lookupd-http-address=127.0.0.1:4161 --http-address=127.0.0.1:4171

# Run a consumer for ticker topic
$ sudo rkt run --net=host --interactive corpix.github.io/nsq:1.0.0 -- nsq_tail --nsqd-tcp-address=127.0.0.1:4150 --topic=ticker

```

#### For Kafka

##### Docker

``` console
$ sudo docker-compose up etcd zetcd kafka

# Run a consumer for ticker topic
$ sudo docker run -it docker.io/wurstmeister/kafka /opt/kafka/bin/kafka-console-consumer.sh -- --from-beginning --zookeeper=127.0.0.1:2181 --topic=ticker
```

##### Rkt

> All commands should be run in separate terminal windows.

``` console
$ sudo rkt run --interactive coreos.com/etcd:v3.1.8 --net=host -- --log-output=stderr --debug
$ sudo rkt run --interactive corpix.github.io/zetcd:0.0.2 --net=host -- --zkaddr=127.0.0.1:2181 --endpoints=127.0.0.1:2379 --logtostderr

# Init zetcd with data(or kafka will fail to start)
$ sudo rkt run                                 \
    --interactive corpix.github.io/zetcd:0.0.2 \
    --net=host --exec=/bin/bash                \
    -- -c "
        zkctl create '/' ''
        zkctl create '/brokers' ''
        zkctl create '/brokers/ids' ''
        zkctl create '/brokers/topics' ''
    "

$ sudo rkt run --interactive corpix.github.io/kafka:2.12-0.10.2.1-1496226351 --net=host

# Run a consumer for ticker topic
$ sudo rkt run --interactive corpix.github.io/kafka:2.12-0.10.2.1-1496226351 --net=host --exec=/usr/bin/kafka-console-consumer -- --from-beginning --zookeeper=127.0.0.1:2181 --topic=ticker
```

#### Running fetcher

Yo will need a fetcher service which will push data into the queue, please see [market-fetcher Preparation section](https://github.com/cryptounicorns/market-fetcher#preparations).

### Running fetcher HTTP API

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
