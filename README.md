# snoopy
HTTP reverse proxy for snooping your services

## Usage

Proxy a given host:

    snoopy --name Example -port 8080 -url http://example.com
    
Show request and response bodies:

    snoopy --name Example -port 8080 -url http://example.com -showBody
    
Configure multiple proxies via file:

    snoopy -file=sample.yml
    
    # sample.yml
    showBody: false
    proxyConfigs:
      -
        name: Example
        port: 8080
        url: http://www.example.com
      -
        name: Google
        port: 9090
        url: http://www.google.com
