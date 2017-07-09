market-fetcher
---------

[![Build Status](https://travis-ci.org/corpix/market-fetcher.svg?branch=master)](https://travis-ci.org/corpix/market-fetcher)

Fetches crypto currency market data and streams it into configured endpoint.

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
go run ./market-fetcher/market-fetcher.go --debug
```
