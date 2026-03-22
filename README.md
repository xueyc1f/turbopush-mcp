# TurboPush MCP Server

TurboPush 的 MCP (Model Context Protocol) Server，让 Claude 等大模型通过标准 MCP 协议直接调用 TurboPush 的内容发布能力。

## 架构

```
Claude (Claude Code / Claude Desktop)
  │
  │  MCP 协议 (stdio)
  │
  ▼
turbo-push-mcp (本项目)
  │
  │  HTTP REST API
  │
  ▼
TurboPush 服务 (127.0.0.1:{port})
```

MCP Server 作为独立进程运行，通过 HTTP 代理方式调用 TurboPush 现有 REST API，对主服务零侵入。

## 编译

需要 Go 1.25+。

```bash
cd mcp
go build -o turbo-push-mcp .
```

国内环境如遇网络问题，可设置代理：

```bash
GOPROXY=https://goproxy.cn,direct go build -o turbo-push-mcp .
```

## 配置

MCP Server 通过环境变量连接 TurboPush 服务：

| 环境变量 | 说明 |
|---------|------|
| `TURBO_PUSH_PORT` | TurboPush 服务端口 |
| `TURBO_PUSH_AUTH` | TurboPush 认证 Token |

TurboPush 每次启动时会生成随机端口和 Token，可从启动日志中获取。

### Claude Code

编辑 `~/.claude/settings.json`（全局）或项目目录下 `.claude/settings.json`：

```json
{
  "mcpServers": {
    "turbo-push": {
      "command": "/绝对路径/mcp/turbo-push-mcp",
      "env": {
        "TURBO_PUSH_PORT": "12345",
        "TURBO_PUSH_AUTH": "你的token..."
      }
    }
  }
}
```

### Claude Desktop

编辑 `~/Library/Application Support/Claude/claude_desktop_config.json`（macOS）：

```json
{
  "mcpServers": {
    "turbo-push": {
      "command": "/绝对路径/mcp/turbo-push-mcp",
      "env": {
        "TURBO_PUSH_PORT": "12345",
        "TURBO_PUSH_AUTH": "你的token..."
      }
    }
  }
}
```

### 手动验证

```bash
TURBO_PUSH_PORT=12345 TURBO_PUSH_AUTH=xxx ./turbo-push-mcp
```

启动后会通过 stdin/stdout 进行 MCP 通信，可用 [MCP Inspector](https://github.com/modelcontextprotocol/inspector) 调试。

## 可用 Tools

共 18 个 Tool，覆盖完整发布流程：

### 平台

| Tool | 说明 |
|------|------|
| `list_platforms` | 获取支持的发布平台列表 |

### 账号

| Tool | 说明 |
|------|------|
| `list_accounts` | 获取所有平台账号 |
| `list_logged_accounts` | 获取已登录的账号 |

### 平台配置

| Tool | 说明 |
|------|------|
| `list_platform_settings` | 获取平台配置列表 |
| `create_platform_setting` | 创建平台配置 |
| `update_platform_setting` | 更新平台配置 |
| `delete_platform_setting` | 删除平台配置 |

### 内容管理

| Tool | 说明 |
|------|------|
| `list_articles` | 获取内容列表 |
| `get_article` | 获取内容详情 |
| `create_article` | 创建文章 |
| `create_graph_text` | 创建图文 |
| `create_video` | 创建视频 |
| `update_article` | 更新内容 |
| `delete_article` | 删除内容 |

### 发布

| Tool | 说明 |
|------|------|
| `publish_article` | 发布文章到指定账号 |
| `publish_graph_text` | 发布图文到指定账号 |
| `publish_video` | 发布视频到指定账号 |

### 发布记录

| Tool | 说明 |
|------|------|
| `list_records` | 获取发布记录列表 |
| `get_record_info` | 获取发布记录详情 |

## 典型工作流

在 Claude 中可以这样使用：

```
> 帮我查看有哪些已登录的抖音账号

> 创建一篇文章，标题"产品更新公告"，内容为 ...

> 把这篇文章发布到所有已登录的微信公众号账号
```

Claude 会自动编排调用：`list_logged_accounts` → `create_article` → `publish_article`。

### 发布参数示例

发布时需要构造 `postAccounts` 数组：

```json
{
  "article_id": 1,
  "postAccounts": [
    {
      "id": 10,
      "platName": "抖音账号A",
      "settings": {
        "platType": "douyin",
        "allowSave": true,
        "lookScope": 0
      }
    }
  ]
}
```

`settings.platType` 对应平台标识：

| platType | 平台 |
|----------|------|
| `wechat` | 微信公众号 |
| `wechat-video` | 微信视频号 |
| `douyin` | 抖音 |
| `toutiaohao` | 今日头条 |
| `kuaishou` | 快手 |
| `xiaohongshu` | 小红书 |
| `bilibili` | 哔哩哔哩 |
| `zhihu` | 知乎 |
| `sina` | 新浪微博 |
| `csdn` | CSDN |
| `juejin` | 掘金 |
| `jianshuhao` | 简书 |
| `tiktok` | TikTok |
| `youtube` | YouTube |
| `x` | X (Twitter) |
| `pinduoduo` | 拼多多 |
| `acfun` | AcFun |
| `omtencent` | 企鹅号 |
| `weishi` | 微视 |
| `baijiahao` | 百家号 |

各平台 `settings` 的完整字段说明见 `docs/api.md` 中的 **setting 参数说明** 章节。

## Skills 集成

### OpenClaw

将 `mcp/skills/turbo-push/` 目录复制到 OpenClaw 的 skills 目录：

```bash
cp -r mcp/skills/turbo-push ~/.openclaw/workspace/skills/
```

重启 OpenClaw 或刷新 skills 即可使用。Skill 会自动配置 MCP Server 连接。

### Claude Code

项目已内置 3 个 slash command（位于 `.claude/commands/`）：

| 命令 | 说明 |
|------|------|
| `/publish` | 发布内容到指定平台 |
| `/publish-all` | 批量发布到所有已登录账号 |
| `/status` | 查看账号和发布状态 |

使用示例：

```
/publish 把这篇文章发到所有抖音账号
/publish-all 标题"新品上线" 内容为...
/status 查看最近的发布记录
```
