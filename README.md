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
