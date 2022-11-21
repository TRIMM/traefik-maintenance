# Traefik Maintenance Middleware Plugin

This Traefik middleware plugin allows you to configure maintenance responses for your routers.  
You have to declare the experimental block in your traefik static configuration file or add the
required flags.

Maintenance mode will be triggered if `enabled` is set to `true` and if the file configured for  
`triggerFilename` exists.

It's also possible to provide a JSON (or any other) maintenance response by changing the
`filename` to point to a JSON file and by changing `httpContentType` to `application/json; charset=utf-8`.

## Static Configuration

### FILE

```yaml
experimental:
  plugins:
    traefik-maintenance:
      moduleName: github.com/TRIMM/traefik-maintenance
      version: v1.0.1
```

### CLI

```shell
--experimental.plugins.traefik-maintenance.modulename=github.com/TRIMM/traefik-maintenance
--experimental.plugins.traefik-maintenance.version=v1.0.1
```

## Dynamic Configuration

### FILE

```yaml
http:
  services:
    service1:
      loadBalancer:
        servers:
          - url: "http://service1:8080/"
    service2:
      loadBalancer:
        servers:
          - url: "http://service2:8081/"
  routers:
    service1-router:
      rule: "Host(`service1`)"
      service: "service1"
      middlewares:
        - maintenance
    service2-router:
      rule: "Host(`service2`)"
      service: "service2"
      middlewares:
        - maintenance
  middlewares:
    maintenance:
      plugin:
        traefik-maintenance:
          enabled: true
          filename: '/path/to/maintenance.html'
          triggerFilename: '/path/to/maintenance.trigger'
          httpResponseCode: 503
          httpContentType: 'text/html; charset=utf-8'
```
