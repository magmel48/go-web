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
DATABASE_DSN="postgres://postgres:Qwerty\!234@localhost:5432/yandex-diploma?sslmode=disable" go test -bench=. -memprofile=profiles/base.pprof
```

```shell
go tool pprof -http=":9090" profiles/base.pprof
```

Result diff for userlinks after changes in package (removing append):
```text
Type: alloc_space
Time: Mar 10, 2022 at 2:27pm (MSK)
Showing nodes accounting for -33.94MB, 1.82% of 1864.95MB total
```

Anyway, the relational database is the most problem.

## Documentation

Available if you run:

```shell
godoc -http=:6060
```

To check documentation from `internal` package here is link (`?m=all` is essential here):

```text
http://localhost:6060/pkg/github.com/magmel48/go-web/internal/app/?m=all
```
