package main

import "strings"

func formatNginxConf(port string) string {
	proxyTarget := "http://host.docker.internal:" + port
	return strings.Replace(NginxConfigTemplate, "{{ PROXY_TARGET }}", proxyTarget, 1)
}

const NginxConfigTemplate = `user  nginx;
worker_processes  auto;

pid  /var/run/nginx.pid;

events {
  worker_connections  1024;
}

http {
  server {
    listen 80;
    location / {
      proxy_pass {{ PROXY_TARGET }};
    }
  }
}
`
