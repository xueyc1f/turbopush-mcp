查看 TurboPush 当前状态。使用 turbo-push MCP 工具完成。

## 步骤

1. 调用 `list_logged_accounts` 获取已登录账号
2. 调用 `list_platforms` 获取支持的平台列表
3. 如用户想查看发布记录，调用 `list_records`
4. 如用户指定了记录ID，调用 `get_record_info` 获取详情

汇总输出：
- 已登录账号数量和列表（按平台分组）
- 最近发布记录（如有）
- 各平台支持的内容类型

$ARGUMENTS
