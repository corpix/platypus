package handshake

/*
What we will get from client:
>=http/1.1
GET request
Upgrade == websocket (case insensitive)
Connection == Upgrade (case insensitive)
Sec-Websocket-Key, 16 bytes in length when decoded(base64.StdEncoding.EncodedLen(x))
Sec-Websocket-Version == 13
Sec-Websocket-Protocols == ?
Sec-Websocket-Extensions == ?
If Origin exists then check match, otherwise ignore
*/

/*

 */
