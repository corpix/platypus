platypus
---------

[![Build Status](https://travis-ci.org/cryptounicorns/platypus.svg?branch=master)](https://travis-ci.org/cryptounicorns/platypus)

Receives data from user configured message queue and broadcasts it to websocket clients.

## Running

### Docker
``` console
$ docker-compose up
```

### No isolation

> From the root of `mole` repository.

``` console
$ go run ./platypus/platypus.go --debug
```

## Building

Build a binary release:

``` console
$ GOOS=linux make
# This will put a binary into ./build/platypus
```


## Testing

> wscp is https://github.com/corpix/wscp

``` console
$ wscp ws://127.0.0.1:8080/api/v1/tickers/stream
{"high":2943.7997,"low":2745.7621,"avg":2844.7809,"vol":1584868.5,"volCur":564.85726,"last":2745.7621,"buy":2766.2631,"sell":2788.4115,"timestamp":1503190041,"currencyPair":"LTC-RUB","market":"yobit"}
{"high":4350,"low":4116.798,"avg":0,"vol":34787.59018931,"volCur":1119.71769938,"last":4280.783,"buy":4280.8723,"sell":4285.3211,"timestamp":1503190091,"currencyPair":"BTC-USD","market":"cex"}
{"high":266882.1,"low":252200.02,"avg":0,"vol":60.3862276,"volCur":2.35581065,"last":266882.1,"buy":256782.46,"sell":263975,"timestamp":1503190091,"currencyPair":"BTC-RUB","market":"cex"}
...
```

## Name origin

> Perry is the pet platypus of the Flynn-Fletcher family, and is perceived as mindless and domesticated. In secret, however, he lives a __double life__ as a member of an all-animal espionage organization referred to as O.W.C.A. (Organization Without a Cool Acronym).

https://en.wikipedia.org/wiki/Perry_the_Platypus

This service lives a double life, kind of :)
