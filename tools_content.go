package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerContentTools(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("list_articles",
			mcp.WithDescription("获取内容列表（文章:publish_type=1/图文:publish_type=2/视频:publish_type=3），支持按状态筛选和分页"),
			mcp.WithNumber("status", mcp.Description("状态筛选：1=草稿，2=已发布")),
			mcp.WithNumber("current", mcp.Description("页码，默认1")),
			mcp.WithNumber("size", mcp.Description("每页条数，默认10")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			path := "/article/list?simple=true&"
			if v, ok := args["status"].(float64); ok {
				path += fmt.Sprintf("status=%d&", int(v))
			}
			if v, ok := args["current"].(float64); ok {
				path += fmt.Sprintf("current=%d&", int(v))
			}
			if v, ok := args["size"].(float64); ok {
				path += fmt.Sprintf("size=%d&", int(v))
			}
			resp, err := c.Get(path)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_article",
			mcp.WithDescription("获取内容详情，返回完整的文章/图文/视频数据"),
			mcp.WithNumber("article_id", mcp.Required(), mcp.Description("内容ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			rid, err := request.RequireFloat("article_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, err := c.Get(fmt.Sprintf("/article/get/%d", int(rid)))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_article",
			mcp.WithDescription("创建文章内容（Markdown 格式），创建后可通过 publish_article 发布到各平台"),
			mcp.WithString("title", mcp.Required(), mcp.Description("文章标题")),
			mcp.WithString("markdown", mcp.Required(), mcp.Description("文章 Markdown 内容")),
			mcp.WithString("desc", mcp.Description("文章摘要/描述")),
			mcp.WithArray("thumb", mcp.Description("封面图路径数组，可传1张或3张，为空则不设置封面")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			body := map[string]any{
				"title":    args["title"],
				"markdown": args["markdown"],
			}
			if v, ok := args["desc"]; ok {
				body["desc"] = v
			}
			if v, ok := args["thumb"]; ok {
				body["thumb"] = v
				body["autoThumb"] = false
			} else {
				body["autoThumb"] = true
			}
			resp, err := c.Post("/article/create", body)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_graph_text",
			mcp.WithDescription("创建图文内容，需要提供图片文件路径。创建后可通过 publish_graph_text 发布到各平台"),
			mcp.WithString("title", mcp.Required(), mcp.Description("图文标题")),
			mcp.WithString("desc", mcp.Description("图文描述，支持使用'#话题名称#'设置话题, 使用'@用户名称 '提及用户，如：#每日发文##解放生产力#@TurboPush @luster 描述内容...")),
			mcp.WithArray("files", mcp.Description("图片文件路径数组")),
			mcp.WithArray("thumb", mcp.Description("封面图路径数组，为空则自动生成")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			body := map[string]any{
				"title": args["title"],
			}
			if v, ok := args["desc"]; ok {
				body["desc"] = v
			}
			if v, ok := args["files"]; ok {
				body["files"] = v
			}
			if v, ok := args["thumb"]; ok {
				body["thumb"] = v
				body["autoThumb"] = false
			} else {
				body["autoThumb"] = true
			}
			resp, err := c.Post("/article/graphText", body)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_video",
			mcp.WithDescription("创建视频内容，需要提供视频文件路径。创建后可通过 publish_video 发布到各平台"),
			mcp.WithString("title", mcp.Required(), mcp.Description("视频标题")),
			mcp.WithString("desc", mcp.Description("视频描述，支持使用'#话题名称#'设置话题, 使用'@用户名称 '提及用户，如：#每日发文##解放生产力#@TurboPush @luster 描述内容...")),
			mcp.WithArray("files", mcp.Required(), mcp.Description("视频文件路径数组")),
			mcp.WithArray("thumb", mcp.Description("封面图路径数组，为空则自动生成")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			body := map[string]any{
				"title": args["title"],
				"files": args["files"],
			}
			if v, ok := args["desc"]; ok {
				body["desc"] = v
			}
			if v, ok := args["thumb"]; ok {
				body["thumb"] = v
				body["autoThumb"] = false
			} else {
				body["autoThumb"] = true
			}
			resp, err := c.Post("/article/video", body)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("update_article",
			mcp.WithDescription("更新已有的内容（文章/图文/视频）"),
			mcp.WithNumber("article_id", mcp.Required(), mcp.Description("内容ID")),
			mcp.WithString("title", mcp.Description("标题")),
			mcp.WithString("markdown", mcp.Description("Markdown 内容")),
			mcp.WithString("desc", mcp.Description("描述")),
			mcp.WithArray("files", mcp.Description("文件路径数组")),
			mcp.WithArray("thumb", mcp.Description("封面图路径数组")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			rid, err := request.RequireFloat("article_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := make(map[string]any)
			for _, key := range []string{"title", "markdown", "desc", "files", "thumb"} {
				if v, ok := args[key]; ok {
					body[key] = v
				}
			}
			resp, err := c.Post(fmt.Sprintf("/article/update/%d", int(rid)), body)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)

	s.AddTool(
		mcp.NewTool("delete_article",
			mcp.WithDescription("删除内容（文章/图文/视频）"),
			mcp.WithNumber("article_id", mcp.Required(), mcp.Description("内容ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			rid, err := request.RequireFloat("article_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, err := c.Delete(fmt.Sprintf("/article/delete/%d", int(rid)))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(resp), nil
		},
	)
}
