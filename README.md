# sisyphus

A tiny load balancer in Go.

## Instructions

```console
# Build sisyphus, summit and zeus
$ make

# Format the entire project
$ make fmt

# Clean the built binaries
$ make clean

# Start sisyphus
$ PORT=8000 ./bin/sisyphus

# Start summit
$ PORT=3000 ./bin/summit

# Test using curl
$ curl localhost:8000/health

# Auto load-testing using zeus. Provide the address of the load-balancer (sisyphus).
$ HOST='localhost:8000' ./bin/zeus
```

## Load Balancer Configuration

```json
# Add a config.json file in /data directory with the following content structure:

{
  "servers": ["<ip1>:<port1>", "<ip2>:<port2>", ...],
  "weights": [<x1>, <x2>, ...],
  "strategy": "<strategy>"
}
```

> Note that \<text> are placeholders. Replace them with the actual values.

```json
# Example config.json file:

{
  "servers": ["127.0.0.1:3000", "127.0.0.1:3001", "127.0.0.1:3002", "127.0.0.1:3003"],
  "weights": [1, 2, 3, 4],
  "strategy": "round-robin"
}
```
