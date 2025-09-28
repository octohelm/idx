fmt:
    go tool gofumpt -w -l .

test:
    CGO_ENABLED=0 go test -count=1 -v -failfast ./...

test-bench:
    CGO_ENABLED=0 go test -count=1  -test.bench=. -test.benchmem ./...

test-race:
    CGO_ENABLED=1 go test -count=1 -race ./...

dep:
    go mod tidy

update:
    go get -u -t ./...
