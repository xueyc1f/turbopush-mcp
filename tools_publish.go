package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPublishTools(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("publish_article",
			mcp.WithDescription(`发布文章到指定平台账号。需要先通过 create_article 创建文章内容，或通过 list_articles 选择已有的文章(publish_type=1)，然后调用此工具发布。
postAccounts 数组中每个元素需要包含：id(账号ID)、platName(平台名)、settings(平台配置，必须包含 platType 字段)。
发布前请先调用 get_platform_setting_schema 查询目标平台的 article 类型所需的 settings 字段。

【定时发布】如需定时发布，请在 settings 中添加 timerPublish 字段：
{
  "timerPublish": {
    "enable": true,
    "timer": "2025-04-25 15:54:00"
  }
}

发布是同步操作，会等待所有账号发布完成后返回结果汇总。`),
			mcp.WithNumber("article_id", mcp.Required(), mcp.Description("文章ID（通过 create_article 创建后获得或通过 list_articles 获取）")),
			mcp.WithBoolean("syncDraft", mcp.Description("是否仅同步草稿（不直接发布），默认 false")),
			mcp.WithBoolean("headless", mcp.Description("是否使用无头浏览器模式，默认 false")),
			mcp.WithArray("postAccounts", mcp.Required(), mcp.Description("发布目标账号数组，每个元素包含 id(账号ID)、platName(平台名)、settings(平台配置对象，含 platType)")),
		),
		publishHandler(c, "/sse/article/%d", ContentArticle),
	)

	s.AddTool(
		mcp.NewTool("publish_graph_text",
			mcp.WithDescription(`发布图文到指定平台账号。需要先通过 create_graph_text 创建图文内容，或通过 list_articles 选择已有的图文(publish_type=2)，然后调用此工具发布。
发布前请先调用 get_platform_setting_schema 查询目标平台的 graph_text 类型所需的 settings 字段。

【定时发布】如需定时发布，请在 settings 中添加 timerPublish 字段：
{
  "timerPublish": {
    "enable": true,
    "timer": "2025-04-25 15:54:00"
  }
}

参数说明同 publish_article。`),
			mcp.WithNumber("article_id", mcp.Required(), mcp.Description("图文ID（通过 create_graph_text 创建后获得或通过 list_articles 获取）")),
			mcp.WithBoolean("syncDraft", mcp.Description("是否仅同步草稿")),
			mcp.WithBoolean("headless", mcp.Description("是否使用无头浏览器模式")),
			mcp.WithArray("postAccounts", mcp.Required(), mcp.Description("发布目标账号数组")),
		),
		publishHandler(c, "/sse/graphText/%d", ContentGraphText),
	)

	s.AddTool(
		mcp.NewTool("publish_video",
			mcp.WithDescription(`发布视频到指定平台账号。需要先通过 create_video 创建视频内容，或通过 list_articles 选择已有的视频(publish_type=3)，然后调用此工具发布。
发布前请先调用 get_platform_setting_schema 查询目标平台的 video 类型所需的 settings 字段。

【定时发布】如需定时发布，请在 settings 中添加 timerPublish 字段：
{
  "timerPublish": {
    "enable": true,
    "timer": "2025-04-25 15:54:00"
  }
}

参数说明同 publish_article。`),
			mcp.WithNumber("article_id", mcp.Required(), mcp.Description("视频ID（通过 create_video 创建后获得或通过 list_articles 获取）")),
			mcp.WithBoolean("syncDraft", mcp.Description("是否仅同步草稿")),
			mcp.WithBoolean("headless", mcp.Description("是否使用无头浏览器模式")),
			mcp.WithArray("postAccounts", mcp.Required(), mcp.Description("发布目标账号数组")),
		),
		publishHandler(c, "/sse/video/%d", ContentVideo),
	)
}

func publishHandler(c *Client, pathTemplate string, contentType string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		rid, err := request.RequireFloat("article_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// 验证并填充 postAccounts 中的 settings
		postAccounts, ok := args["postAccounts"].([]any)
		if !ok || len(postAccounts) == 0 {
			return mcp.NewToolResultError("postAccounts 不能为空"), nil
		}

		if errMsg := validateAndFillDefaults(postAccounts, contentType); errMsg != "" {
			return mcp.NewToolResultError(errMsg), nil
		}

		body := map[string]any{
			"postAccounts": postAccounts,
		}
		if v, ok := args["syncDraft"].(bool); ok {
			body["syncDraft"] = v
		}
		if v, ok := args["headless"].(bool); ok {
			body["headless"] = v
		}

		path := fmt.Sprintf(pathTemplate, int(rid))
		events, err := c.PostSSE(path, body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("发布请求失败: %s", err)), nil
		}

		return mcp.NewToolResultText(summarizeEvents(events)), nil
	}
}

// validateAndFillDefaults 校验 postAccounts 中每个账号的 settings 必填字段，
// 并自动填充 schema 中有默认值但 settings 中未设置的字段。
// 返回空字符串表示校验通过，否则返回错误消息。
func validateAndFillDefaults(postAccounts []any, contentType string) string {
	var errs []string

	for i, acc := range postAccounts {
		accMap, ok := acc.(map[string]any)
		if !ok {
			errs = append(errs, fmt.Sprintf("postAccounts[%d]: 格式无效，需要对象", i))
			continue
		}

		settings, ok := accMap["settings"].(map[string]any)
		if !ok {
			errs = append(errs, fmt.Sprintf("postAccounts[%d]: 缺少 settings 对象", i))
			continue
		}

		platType, _ := settings["platType"].(string)
		if platType == "" {
			errs = append(errs, fmt.Sprintf("postAccounts[%d]: settings 中缺少 platType 字段", i))
			continue
		}

		// 获取 schema
		fields, ok := getSchema(platType, contentType)
		if !ok {
			// 平台没有 schema 定义，跳过验证（允许通过，由后端处理）
			continue
		}

		// 检查必填字段
		var missing []string
		for _, f := range fields {
			if f.Required {
				if _, exists := settings[f.Name]; !exists {
					missing = append(missing, f.Name)
				}
			}
		}
		if len(missing) > 0 {
			platName, _ := accMap["platName"].(string)
			if platName == "" {
				platName = platType
			}
			errs = append(errs, fmt.Sprintf(
				"postAccounts[%d] (%s): 缺少必填字段 [%s]，请先调用 get_platform_setting_schema 查看所需字段",
				i, platName, strings.Join(missing, ", ")))
		}

		// 自动填充默认值
		for _, f := range fields {
			if f.Default != nil {
				if _, exists := settings[f.Name]; !exists {
					settings[f.Name] = f.Default
				}
			}
		}
	}

	if len(errs) > 0 {
		return "发布参数校验失败:\n" + strings.Join(errs, "\n")
	}
	return ""
}

func summarizeEvents(events []sseEvent) string {
	var (
		successes []string
		errors    []string
		result    string
	)

	for _, ev := range events {
		switch ev.Event {
		case "success":
			successes = append(successes, ev.Data)
		case "error":
			errors = append(errors, ev.Data)
		case "finish":
			var finish struct {
				Msg string `json:"msg"`
				Res []bool `json:"res"`
			}
			if json.Unmarshal([]byte(ev.Data), &finish) == nil {
				total := len(finish.Res)
				successCount := 0
				for _, r := range finish.Res {
					if r {
						successCount++
					}
				}
				result = fmt.Sprintf("%s (成功: %d/%d)", finish.Msg, successCount, total)
			} else {
				result = ev.Data
			}
		case "wait":
			return fmt.Sprintf("发布等待: %s", ev.Data)
		case "vip":
			return fmt.Sprintf("会员限制: %s", ev.Data)
		}
	}

	var sb strings.Builder
	if result != "" {
		sb.WriteString("## 发布结果\n")
		sb.WriteString(result)
		sb.WriteString("\n\n")
	}
	if len(successes) > 0 {
		sb.WriteString("### 成功\n")
		for _, s := range successes {
			sb.WriteString("- ")
			sb.WriteString(s)
			sb.WriteString("\n")
		}
	}
	if len(errors) > 0 {
		sb.WriteString("\n### 错误\n")
		for _, e := range errors {
			sb.WriteString("- ")
			sb.WriteString(e)
			sb.WriteString("\n")
		}
	}
	if sb.Len() == 0 {
		return "发布完成，无详细事件信息"
	}
	return sb.String()
}
