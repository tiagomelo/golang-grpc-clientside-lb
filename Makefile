# Copyright (c) 2023 Tiago Melo. All rights reserved.
# Use of this source code is governed by the MIT License that can be found in
# the LICENSE file.

include .env
export

# ==============================================================================
# Help

.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]\n"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==============================================================================
# Protofile compilation

.PHONY: proto
## proto: compile proto files
proto:
	@ rm -rf api/proto/gen/helloservice
	@ mkdir -p api/proto/gen/helloservice
	@ cd api/proto ; \
	protoc --go_out=gen/helloservice --go_opt=paths=source_relative --go-grpc_out=gen/helloservice --go-grpc_opt=paths=source_relative helloservice.proto

# ==============================================================================
# Loopback addresses

.PHONY: loopback-up
## loopback-up: adds loopback addresses
loopback-up:
	@ sudo ifconfig lo0 alias $(LOOPBACK_ADDRESS_1) && sudo ifconfig lo0 alias $(LOOPBACK_ADDRESS_2)

.PHONY: loopback-down
## loopback-down: remove loopback addresses
loopback-down:
	@ sudo ifconfig lo0 -alias $(LOOPBACK_ADDRESS_1) && sudo ifconfig lo0 -alias $(LOOPBACK_ADDRESS_2)

.PHONY: show-loopback-addrs
## show-loopback-addrs: show loopback addresses
show-loopback-addrs:
	@ ifconfig lo0

# ==============================================================================
# CoreDNS

.PHONY: build-coredns-img
## build-coredns-img: builds CoreDNS docker image
build-coredns-img:
	docker build --no-cache -t grpc-coredns -f ./config/coredns/Dockerfile .

.PHONY: coredns
## coredns: runs CoreDNS
coredns: build-coredns-img
	@ docker run --name grpc-coredns -p 53:53/udp --rm grpc-coredns

# ==============================================================================
# socat TCP proxy

.PHONY: socat-server-one
## socat-server-one: starts socat TCP proxy for server 1
socat-server-one:
	@ socat TCP-LISTEN:$(PROM_TARGET_GRPC_SERVER_ONE_PORT),fork TCP:$(LOOPBACK_ADDRESS_1):$(SOCAT_GRPC_SERVER_ONE_PORT)
	@ echo "socat is running for server 1... hit control + c to stop it."

.PHONY: socat-server-two
## socat-server-two: starts socat TCP proxy for server 2
socat-server-two:
	@ socat TCP-LISTEN:$(PROM_TARGET_GRPC_SERVER_TWO_PORT),fork TCP:$(LOOPBACK_ADDRESS_2):$(SOCAT_GRPC_SERVER_TWO_PORT)
	@ echo "socat is running for server 2... hit control + c to stop it."

# ==============================================================================
# Metrics

.PHONY: parse-templates
## parse-templates: parses Prometheus scrapes and datasource templates
parse-templates:
	@ go run templateparser/templateparser.go

.PHONY: obs
## obs: runs both prometheus and grafana
obs: parse-templates
	@ docker-compose up

.PHONY: obs-stop
## obs-stop: stops both prometheus and grafana
obs-stop:
	@ docker-compose down -v

# ==============================================================================
# gRPC server execution

.PHONY: server
## server: runs gRPC server at specified port
server:
	@ if [ -z "$(SERVER_HOST)" ]; then echo >&2 please set the server host via the variable SERVER_HOST; exit 2; fi
	@ if [ -z "$(SERVER_PORT)" ]; then echo >&2 please set the server port via the variable SERVER_PORT; exit 2; fi
	@ if [ -z "$(METRICS_SERVER_PORT)" ]; then echo >&2 please set the metrics server port via the variable METRICS_SERVER_PORT; exit 2; fi
	@ go run cmd/main.go -s $(SERVER_HOST) -p $(SERVER_PORT) -x $(METRICS_SERVER_PORT)

# ==============================================================================
# gRPC client execution

.PHONY: client-round-robin
## client-round-robin: runs gRPC client with round-robin load balancing policy
client-round-robin:
	@ go run client/client.go -l round_robin

.PHONY: client-pick-first
## client-pick-first: runs gRPC client with pick-first load balancing policy
client-pick-first:
	@ go run client/client.go -l pick_first

