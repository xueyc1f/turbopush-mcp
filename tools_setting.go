package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSettingTools(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("list_platform_settings",
			mcp.WithDescription("获取指定平台的配置列表，每个配置包含平台特定的发布参数（如定时发布、原创声明、可见范围等）"),
			mcp.WithNumber("platform_id", mcp.Required(), mcp.Description("平台ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			pid, err := request.RequireFloat("platform_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, err := c.Get(fmt.Sprintf("/platSet/list/%d", int(pid)))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_platform_setting",
			mcp.WithDescription(`创建平台发布配置。setting 对象必须包含 platType 字段，其余字段根据平台不同而不同。
常见平台 platType：wechat(微信公众号)、wechat-video(视频号)、douyin(抖音)、toutiaohao(头条)、kuaishou(快手)、xiaohongshu(小红书)、bilibili(B站)、zhihu(知乎)、sina(微博)、csdn、juejin(掘金)、tiktok、youtube、x(Twitter)、pinduoduo(拼多多)等`),
			mcp.WithString("name", mcp.Required(), mcp.Description("配置名称")),
			mcp.WithString("description", mcp.Description("配置描述")),
			mcp.WithNumber("platform_id", mcp.Required(), mcp.Description("平台ID")),
			mcp.WithObject("setting", mcp.Required(), mcp.Description("平台配置内容，JSON对象，必须包含 platType 字段")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			pid, err := request.RequireFloat("platform_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := map[string]any{
				"name":        args["name"],
				"platform_id": int(pid),
				"setting":     args["setting"],
			}
			if desc, ok := args["description"]; ok {
				body["description"] = desc
			}
			resp, err := c.Post("/platSet/create", body)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("update_platform_setting",
			mcp.WithDescription("更新已有的平台发布配置"),
			mcp.WithNumber("setting_id", mcp.Required(), mcp.Description("配置ID")),
			mcp.WithString("name", mcp.Description("配置名称")),
			mcp.WithString("description", mcp.Description("配置描述")),
			mcp.WithNumber("platform_id", mcp.Description("平台ID")),
			mcp.WithObject("setting", mcp.Description("平台配置内容")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			sid, err := request.RequireFloat("setting_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := make(map[string]any)
			for _, key := range []string{"name", "description", "platform_id", "setting"} {
				if v, ok := args[key]; ok {
					body[key] = v
				}
			}
			resp, err := c.Post(fmt.Sprintf("/platSet/update/%d", int(sid)), body)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("delete_platform_setting",
			mcp.WithDescription("删除平台发布配置"),
			mcp.WithNumber("setting_id", mcp.Required(), mcp.Description("配置ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sid, err := request.RequireFloat("setting_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, err := c.Delete(fmt.Sprintf("/platSet/delete/%d", int(sid)))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)
}
