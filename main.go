package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mark3labs/mcp-go/server"
)

// mcpConfig 对应 turbo_push 写入的 ~/.TurboPush/mcp.json
type mcpConfig struct {
	Port json.Number `json:"port"`
	Auth string      `json:"auth"`
}

// loadConfig 优先读环境变量，fallback 读 ~/.TurboPush/mcp.json
func loadConfig() (port, auth string, err error) {
	port = os.Getenv("TURBO_PUSH_PORT")
	auth = os.Getenv("TURBO_PUSH_AUTH")
	if port != "" && auth != "" {
		return
	}

	home, e := os.UserHomeDir()
	if e != nil {
		err = fmt.Errorf("无法获取用户目录: %w", e)
		return
	}
	data, e := os.ReadFile(filepath.Join(home, ".TurboPush", "mcp.json"))
	if e != nil {
		err = fmt.Errorf("未找到 ~/.TurboPush/mcp.json，且环境变量 TURBO_PUSH_PORT / TURBO_PUSH_AUTH 未设置: %w", e)
		return
	}
	var cfg mcpConfig
	if e := json.Unmarshal(data, &cfg); e != nil {
		err = fmt.Errorf("解析 mcp.json 失败: %w", e)
		return
	}
	if port == "" {
		port = cfg.Port.String()
		// 确保是有效端口号
		if _, e := strconv.Atoi(port); e != nil {
			err = fmt.Errorf("mcp.json 中 port 无效: %s", port)
			return
		}
	}
	if auth == "" {
		auth = cfg.Auth
	}
	if port == "" || auth == "" {
		err = fmt.Errorf("port 或 auth 为空")
	}
	return
}

func main() {
	port, auth, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	client := NewClient(port, auth)

	s := server.NewMCPServer(
		"turbo-push",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithInstructions(`TurboPush 多平台内容发布工具。

典型工作流程：
1. list_platforms - 查看支持的平台
2. list_logged_accounts - 查看已登录的账号
3. create_article / create_graph_text / create_video - 创建内容
4. publish_article / publish_graph_text / publish_video - 发布到指定账号
5. list_records / get_record_info - 查看发布结果

发布时需要构造 postAccounts 数组，每个元素包含：
- id: 账号ID（从 list_logged_accounts 获取）
- platName: 平台名称
- settings: 平台配置对象，必须包含 platType 字段（如 "wechat"、"douyin" 等）

可通过 list_platform_settings 查看已有配置，或通过 create_platform_setting 创建新配置。`),
	)

	registerTools(s, client)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "MCP Server 错误: %s\n", err)
		os.Exit(1)
	}
}
