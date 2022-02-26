# go-web

fascinating one.

## Pprof

```shell
go test -bench=. -mempofile=profiles/base.pprof
```

```shell
go tool pprof -http=":9090" profiles/base.pprof
```
