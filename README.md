# golang-grpc-clientside-lb

A sample project showing how gRPC handles load balancing at client side.

## prerequisites

- [Docker](docker.com)
- [docker-compose](docs.docker.com/compose/)
- [socat](http://www.dest-unreach.org/socat/)

## running it

These steps needs to be done in the presented order. They work on macOS.

### add loopback addresses

```
make loopback-up
```

double check they were created:

```
make show-loopback-addrs
```

output should include:

```
	inet 127.0.0.2 netmask 0xff000000 
	inet 127.0.0.3 netmask 0xff000000 
```

### start Grafana and Prometheus

```
make obs
```

Then head out to `localhost:3000`. Password is defined in `GF_SECURITY_ADMIN_PASSWORD` env var in `.env` file.

The dashboard name is `ServedGRPCRequestsPerServer`.

To stop Grafana:

```
make obs-stop
```

This will delete both Grafana and Prometheus volumes as well.

### start servers

```
make server SERVER_HOST=127.0.0.2 SERVER_PORT=50051 METRICS_SERVER_PORT=2112
```

```
make server SERVER_HOST=127.0.0.3 SERVER_PORT=50051 METRICS_SERVER_PORT=2113
```

### start socat for server one

```
make socat-server-one
```

### start socat for server two

```
make socat-server-two
```

### start coreDNS

```
make coredns
```

### run the client

For client with round-robin load balancing policy:

```
make client-round-robin
```

For client with pick-first load balancing policy:

```
make client-pick-first
```

## Makefile

```
make help

Usage: make [target]

  help                  shows this help message
  proto                 compile proto files
  loopback-up           adds loopback addresses
  loopback-down         remove loopback addresses
  show-loopback-addrs   show loopback addresses
  build-coredns-img     builds CoreDNS docker image
  coredns               runs CoreDNS
  socat-server-one      starts socat TCP proxy for server 1
  socat-server-two      starts socat TCP proxy for server 2
  parse-templates       parses Prometheus scrapes and datasource templates
  obs                   runs both prometheus and grafana
  obs-stop              stops both prometheus and grafana
  server                runs gRPC server at specified port
  client-round-robin    runs gRPC client with round-robin load balancing policy
  client-pick-first     runs gRPC client with pick-first load balancing policy
```