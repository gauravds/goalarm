# goalarm

Watch the linux performance and see the alarm if threshold crossed.

## Run by docker images:

1. `docker run -it goalarm:1.0`
1. `vmstat 1 | ./main` or `vmstat 1 | go run main.go`

## Command line arguments available as:

`go run main.go -h` or `./main -h`

```sh
-c int
      max threshold count to show panic (default 3)
-t int
      max threshold (default 194700)
-w string
      watch field (default "free")
```

### Use as want to set `cache` alarm with threshold `2405736` and max attempt: 2

`-w=cache -t=2405736 -c=2`

#### as full command:

`vmstat 1 | ./main -w=cache -t=2405736 -c=2`

or

`vmstat 1 | go run main.go -w=cache -t=2405736 -c=2`

## Build docker image locally:

`docker build . --no-cache -t goalarm:1.0`
