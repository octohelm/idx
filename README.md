# idx

`idx` 是一个小型 Go 库仓库，提供与分布式 ID 生成相关的基础能力。

当前仓库聚焦两类职责：`pkg/snowflake` 提供雪花 ID 生成能力，`pkg/workerid` 提供基于 IPv4 地址推导 worker ID 的辅助逻辑。仓库本身不负责服务运行、部署编排或示例应用。

## 入口

- [`pkg/snowflake`](/Users/morlay/src/github.com/octohelm/idx/pkg/snowflake/snowflake.go): 雪花 ID 工厂与生成器实现。
- [`pkg/workerid`](/Users/morlay/src/github.com/octohelm/idx/pkg/workerid/ip.go): 从 IPv4 地址提取 worker ID 的辅助逻辑。
- [`justfile`](/Users/morlay/src/github.com/octohelm/idx/justfile): 统一执行入口，包含依赖整理、格式化和测试命令。
- [`AGENTS.md`](/Users/morlay/src/github.com/octohelm/idx/AGENTS.md): 当前仓库对 agent 生效的协作约束。
