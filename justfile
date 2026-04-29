# 列出仓库可用入口。
[group('meta')]
default:
    @just --list --unsorted

# 整理模块依赖。
[group('dep')]
[no-cd]
dep:
    go mod tidy

# 运行全部测试。
[group('test')]
[no-cd]
test path="./..." *args:
    CGO_ENABLED=0 go test {{ args }} -count=1 -failfast {{ path }}

# 格式化 Go 代码。
[group('dev')]
[no-cd]
fmt:
    go tool gofumpt -w -l .

# 更新直接依赖与测试依赖。
[group('dev')]
[no-cd]
update path='./...':
    go get -u -t {{ path }}
