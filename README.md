# skynet

install:

```
$ go install github.com/btwiuse/skynet/cmd/skynet@latest
```

server:

```
$ skynet server
```

example clients:

```
$ skynet client
$ skynet gos
$ skynet echo
```

TODO

- [ ] reverseproxy WebTransport requests
- [ ] support user specified hostname, requiring netrc authentication
- [x] support custom root domain, for example `HOST=usesthis.app`
- [x] support concurrent rw on session manager map
