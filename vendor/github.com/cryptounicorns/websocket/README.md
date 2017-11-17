websocket
---------

[![Build Status](https://travis-ci.org/cryptounicorns/websocket.svg?branch=master)](https://travis-ci.org/cryptounicorns/websocket)

This is a work-in-progress experiment on websocket library implementation for golang.

It is intended to be as high-level as possible to provide a user useful abstractions
for communication with websocket servers and clients. At the same time it tries to be
as fast as possible and don't restricts you to use low-level things
(such as sending raw websocket frames).

It tries to be compatible with RFC as much as possible.
This project is inspired by https://github.com/gobwas/ws

# Status

Things to implement:

- [ ] consumer
- [ ] producer
- [ ] (describe others)
