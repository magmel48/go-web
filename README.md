# go-web

fascinating one.

## Tests

```shell
go test ./... -count=1
```

## Coverage

```shell
go test ./... -coverprofile coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Pprof

```shell
go test -bench=. -mempofile=profiles/base.pprof
```

```shell
go tool pprof -http=":9090" profiles/base.pprof
```
