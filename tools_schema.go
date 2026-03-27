package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSchemaTools(s *server.MCPServer, _ *Client) {
	s.AddTool(
		mcp.NewTool("get_platform_setting_schema",
			mcp.WithDescription(`查询指定平台在特定内容类型下的 setting 字段定义。
返回每个字段的名称、类型、是否必填、描述、默认值和可选枚举值。
发布内容前务必调用此工具，确保 postAccounts 中的 settings 包含所有必填字段。
不同平台、不同内容类型（文章/图文/视频）的字段可能不同。

【重要】如需定时发布，请检查返回的 schema 中是否包含 timerPublish 字段，如有则按其格式设置。`),
			mcp.WithString("plat_type", mcp.Required(), mcp.Description("平台类型标识，如 wechat、douyin、bilibili、juejin 等")),
			mcp.WithString("content_type", mcp.Required(), mcp.Description("内容类型：article（文章）、graph_text（图文）、video（视频）")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			platType, _ := args["plat_type"].(string)
			contentType, _ := args["content_type"].(string)

			if platType == "" {
				return mcp.NewToolResultError("plat_type 不能为空"), nil
			}
			if contentType == "" {
				return mcp.NewToolResultError("content_type 不能为空，可选值：article, graph_text, video"), nil
			}

			// 校验 content_type
			switch contentType {
			case ContentArticle, ContentGraphText, ContentVideo:
				// ok
			default:
				return mcp.NewToolResultError(fmt.Sprintf(
					"无效的 content_type: %q，可选值：article, graph_text, video", contentType)), nil
			}

			fields, ok := getSchema(platType, contentType)
			if !ok {
				supported := getSupportedPlatTypes()
				return mcp.NewToolResultError(fmt.Sprintf(
					"未找到平台 %q 的 %s 配置 schema。\n支持的平台：%s",
					platType, contentType, strings.Join(supported, ", "))), nil
			}

			data, err := json.MarshalIndent(fields, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("序列化 schema 失败: %s", err)), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
