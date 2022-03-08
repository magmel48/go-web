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

## Benchmarks

```shell
DATABASE_DSN="postgres://postgres:Qwerty\!234@localhost:5432/yandex-diploma?sslmode=disable" go test ./... -bench=.
```

## Pprof

Run per each package separately:
```shell
DATABASE_DSN="postgres://postgres:Qwerty\!234@localhost:5432/yandex-diploma?sslmode=disable" go test -bench=. -mempofile=profiles/base.pprof
```

```shell
go tool pprof -http=":9090" profiles/base.pprof
```
