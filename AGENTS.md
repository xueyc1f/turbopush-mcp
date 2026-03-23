# AGENTS.md - TurboPush MCP Server

本文档为 AI 编码代理提供项目上下文和编码规范。

## 项目概述

TurboPush MCP Server 是一个 Go 语言编写的 MCP (Model Context Protocol) 服务器，作为 Claude 等 AI 模型与 TurboPush 内容发布系统之间的桥梁，支持将内容一键发布到 20+ 平台。

**技术栈**: Go 1.25+, MCP-Go SDK (github.com/mark3labs/mcp-go)

---

## 构建、测试与开发命令

### 构建

```bash
# 标准构建
go build -o turbo-push-mcp .

# 国内环境（使用代理）
GOPROXY=https://goproxy.cn,direct go build -o turbo-push-mcp .

# 生产构建（优化体积）
go build -ldflags="-s -w" -o turbo-push-mcp .

# 交叉编译（示例：Linux ARM64）
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o turbo-push-mcp .
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行单个测试文件
go test -v ./path/to/package -run TestFunctionName

# 运行带覆盖率
go test -cover ./...

# 测试并显示详细输出
go test -v ./...
```

### 代码质量

```bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 检查代码风格（如已安装 golangci-lint）
golangci-lint run
```

### 依赖管理

```bash
# 下载依赖
go mod download

# 整理依赖
go mod tidy

# 验证依赖
go mod verify
```

---

## 项目结构

```
turbopush-mcp/
├── main.go              # 入口：配置加载、MCP Server 初始化
├── client.go            # HTTP Client：API 调用、SSE 流处理
├── schema.go            # 平台 Schema 定义：各平台 settings 字段规范
├── tools.go             # Tool 注册入口
├── tools_*.go           # 各类 MCP Tool 实现
│   ├── tools_platform.go    # 平台相关
│   ├── tools_account.go     # 账号相关
│   ├── tools_setting.go     # 配置相关
│   ├── tools_content.go     # 内容管理
│   ├── tools_publish.go     # 发布操作
│   ├── tools_record.go      # 发布记录
│   └── tools_schema.go      # Schema 查询
├── skills/turbo-push/   # OpenClaw/Claude Code Skill 定义
├── .claude/commands/    # Claude Code Slash Commands
│   ├── publish.md
│   ├── publish-all.md
│   └── status.md
└── docs/api.md          # TurboPush API 完整文档
```

---

## Go 代码风格指南

### 导入规范

```go
// 标准 Go 导入分组：
// 1. 标准库
// 2. 第三方库
// 3. 本地包（本项目无子包）

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)
```

### 命名约定

- **包名**: 小写单词，不使用下划线 (`package main`)
- **导出函数/类型**: PascalCase（如 `NewClient`, `FieldDef`）
- **私有函数/字段**: camelCase（如 `loadConfig`, `do`）
- **常量**: PascalCase 或全大写（如 `ContentArticle`, `ContentAll`）
- **接口**: 通常以 `-er` 结尾（本项目暂无自定义接口）

### 结构体定义

```go
// 字段对齐，json tag 使用蛇形命名
type FieldDef struct {
    Name        string       `json:"name"`
    Type        string       `json:"type"`        // "string", "bool", "uint", "int", "object", "array"
    Required    bool         `json:"required"`
    Description string       `json:"description"`
    Default     any          `json:"default,omitempty"`
    Options     []OptionDef  `json:"options,omitempty"`
}
```

### 错误处理

```go
// 使用 fmt.Errorf 包装错误，保留上下文
if err != nil {
    return nil, fmt.Errorf("create request: %w", err)
}

// MCP Tool 错误返回
if err != nil {
    return mcp.NewToolResultError(err.Error()), nil
}

// 错误消息使用中文（匹配用户群体）
err = fmt.Errorf("未找到 ~/.TurboPush/mcp.json，且环境变量 TURBO_PUSH_PORT / TURBO_PUSH_AUTH 未设置: %w", e)
```

### MCP Tool 注册模式

```go
func registerXxxTools(s *server.MCPServer, c *Client) {
    s.AddTool(
        mcp.NewTool("tool_name",
            mcp.WithDescription("工具描述（中文）"),
            mcp.WithString("param", mcp.Required(), mcp.Description("参数描述")),
        ),
        func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
            // 1. 获取参数
            args := request.GetArguments()
            
            // 2. 参数验证
            rid, err := request.RequireFloat("article_id")
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            
            // 3. 调用 API
            resp, err := c.Get(fmt.Sprintf("/path/%d", int(rid)))
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            
            // 4. 返回结果
            return jsonResult(resp), nil
        },
    )
}
```

### 注释规范

- 导出函数必须有注释，以函数名开头
- 注释使用中文描述

```go
// loadConfig 优先读环境变量，fallback 读 ~/.TurboPush/mcp.json
func loadConfig() (port, auth string, err error) { ... }

// getSchema 查询指定平台和内容类型的 setting schema
// 优先精确匹配 contentType，未命中则 fallback 到 "*"
func getSchema(platType, contentType string) ([]FieldDef, bool) { ... }
```

---

## 重要注意事项

### 环境变量

- `TURBO_PUSH_PORT`: TurboPush 服务端口
- `TURBO_PUSH_AUTH`: 认证 Token
- `GOPROXY`: 国内环境设置为 `https://goproxy.cn,direct`

### API 规范

- Base URL: `http://127.0.0.1:{port}`
- 认证方式: 请求头 `Authorization: {token}`
- 响应格式: `{ "code": 200, "msg": "ok", "data": {} }`

### 平台标识 (platType)

`wechat` | `wechat-video` | `douyin` | `toutiaohao` | `kuaishou` | `xiaohongshu` | `bilibili` | `zhihu` | `sina` | `csdn` | `juejin` | `jianshuhao` | `tiktok` | `youtube` | `x` | `pinduoduo` | `acfun` | `omtencent` | `weishi` | `baijiahao`

### 发布流程

1. `list_logged_accounts` → 获取已登录账号
2. `create_article`/`create_graph_text`/`create_video` → 创建内容
3. `get_platform_setting_schema` → 查询目标平台配置要求
4. 构造 `postAccounts` 数组
5. `publish_*` → 发布到指定账号

---

## 参考文档

- `docs/api.md` - TurboPush 完整 API 文档
- `README.md` - 项目说明和使用指南
- `skills/turbo-push/SKILL.md` - Skill 集成说明
