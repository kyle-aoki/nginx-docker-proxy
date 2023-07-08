### Example Usage

```
# proxies http://0.0.0.0:80 -> http://0.0.0.0:8080
ndp 8080

# proxies http://0.0.0.0:80 -> http://0.0.0.0:8081
ndp 8081
```

### Installation

```
# requires go version 1.20+

go install github.com/kyle-aoki/nginx-docker-proxy
```

### Guide

This docker image will use nginx to reverse proxy requests
from http://0.0.0.0:80 (your host machine's port `80`)
to a port of your choice. There is no down time when switching ports.
This is because nginx can reload its configuration without dropping
any requests.

The port it will proxy to is determined by your configuration
of nginx-docker-proxy. To configure target ports, use
`ndp <port>`.

### How it works

nginx-docker-proxy
