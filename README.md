market-fetcher-http
---------

[![Build Status](https://travis-ci.org/cryptounicorns/market-fetcher-http.svg?branch=master)](https://travis-ci.org/cryptounicorns/market-fetcher-http)

HTTP interface for data providen by [market-fetcher](https://github.com/corpix/market-fetcher).

## Development

> All commands should be run in separate terminal windows.

``` bash
sudo rkt run --net=host --interactive corpix.github.io/nsq:1.0.0 -- nsqd
sudo rkt run --net=host --interactive corpix.github.io/nsq:1.0.0 -- nsqlookupd

# Optionally you could run nsqadmin which will provide you a WEBUI for nsq topics etc.
sudo rkt run --net=host --interactive corpix.github.io/nsq:1.0.0 -- nsqadmin -lookupd-http-address 127.0.0.1:4161

# Run a consumer for ticker topic
sudo rkt run --net=host --interactive corpix.github.io/nsq:1.0.0 -- nsq_tail --nsqd-tcp-address 127.0.0.1:4150 --topic ticker
```

At this point you should be ready to run an application:

``` bash
go run ./market-fetcher-http/market-fetcher-http.go --debug
```

## Testing

> wsd is https://github.com/alexanderGugel/wsd

``` bash
wsd -url ws://127.0.0.1:8080/api/v1/tickers/stream
# < {"High":45.000001,"Low":41.35,"Avg":43.175001,"Vol":13157.821,"VolCur":302.6894,"Last":44.7,"Buy":44.749501,"Sell":44.978,"Timestamp":1501939951,"CurrencyPair":"LTC-USD","Market":"yobit"}
#< {"High":2720,"Low":2504,"Avg":2612,"Vol":1608422.4,"VolCur":616.98825,"Last":2720,"Buy":2670,"Sell":2720,"Timestamp":1501939978,"CurrencyPair":"LTC-RUB","Market":"yobit"}
# ...
```
