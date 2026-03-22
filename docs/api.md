# TurboPush API 文档

## 基础信息

- **Base URL**: `http://127.0.0.1:{port}`
- **认证方式**: 请求头 `Authorization: {token}` 或查询参数 `?auth={token}`
- **登录校验**: 除 `user/login` 外，所有接口需要用户已登录

## 统一响应格式

```json
{
  "code": 200,
  "msg": "ok",
  "data": {}
}
```

| code | 说明 |
|------|------|
| 200  | 成功 |
| 400  | 参数错误 |
| 404  | 未找到 |
| 425  | 非会员限制 |
| 500  | 服务错误 |

## 分页参数

部分列表接口支持分页，查询参数：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| current | int | 1 | 当前页码 |
| size | int | 10 | 每页条数 |

分页响应：
```json
{
  "list": [],
  "pager": {
    "total": 100,
    "size": 10,
    "current": 1
  }
}
```

---

## 系统

### 健康检查

`GET /ping`

无需认证，返回 `200` 状态码。

### 登录状态检查

`GET /check`

返回 `200` 表示已登录且认证通过。

---

## 用户 /user

### 登录

`POST /user/login`

无需登录校验。

**请求体**：
```json
{ "code": "登录授权码" }
```

**响应 data**：用户信息对象

### 登出

`POST /user/logout`

**响应 data**：`null`

### 使用邀请码

`POST /user/invite`

**请求体**：
```json
{ "inviteCode": "邀请码" }
```

**响应 data**：用户信息对象

### 设置头像

`POST /user/avatar`

**请求体**：`multipart/form-data`，字段 `file`（图片文件）

**响应 data**：Base64 头像字符串

### 设置昵称

`POST /user/name`

**请求体**：
```json
{ "name": "昵称" }
```

**响应 data**：昵称字符串

### 保存设置

`POST /user/set`

**请求体**：
```json
{
  "watchInterval": 60,
  "browserType": 1,
  "proxyCheckAddr": "https://www.google.com",
  "fingerBrowser": ""
}
```

| 字段 | 说明 |
|------|------|
| watchInterval | 监听间隔（必填，不为0） |
| browserType | 浏览器类型：1=Chrome, 2=Edge, 3=指纹浏览器 |
| proxyCheckAddr | 代理检测地址（必填） |
| fingerBrowser | 指纹浏览器路径（VIP功能） |

**响应 data**：设置对象

### 获取用户信息

`GET /user/info`

**响应 data**：用户信息对象

---

## 账号 /account

### 登录平台账号

`POST /account/login/:pid`

| 路径参数 | 说明 |
|---------|------|
| pid | 平台ID |

**请求体**：
```json
{
  "proxy": "代理地址（可选，VIP功能）",
  "aid": 0
}
```

**响应 data**：`null`

### 打开账号管理

`POST /account/openManager/:aid`

| 路径参数 | 说明 |
|---------|------|
| aid | 账号ID |

**响应 data**：`null`

### 获取账号列表

`GET /account/list`

**响应 data**：账号数组

### 获取已登录账号

`GET /account/logged`

**响应 data**：已登录账号数组

### 删除账号

`DELETE /account/delete/:aid`

| 路径参数 | 说明 |
|---------|------|
| aid | 账号ID |

**响应 data**：被删除的账号对象

### 一键绑定代理

`POST /account/bind/proxies`

VIP功能。自动将代理分配给账号。

**响应 data**：`null`

### 解除所有代理绑定

`POST /account/unbind/proxies`

**响应 data**：`null`

### 绑定指定代理

`POST /account/bind/proxy`

VIP功能。

**请求体**：
```json
{
  "account_id": 1,
  "proxy_id": 2
}
```

**响应 data**：`null`

### 解除指定账号代理

`POST /account/unbind/proxy/:aid`

| 路径参数 | 说明 |
|---------|------|
| aid | 账号ID |

**响应 data**：`null`

### 设置账号排序

`POST /account/setSort/:aid`

| 路径参数 | 说明 |
|---------|------|
| aid | 账号ID |

**请求体**：
```json
{ "sort": 1 }
```

**响应 data**：账号对象

### 设置账号指纹

`POST /account/setFingerprint/:aid`

VIP功能。

| 路径参数 | 说明 |
|---------|------|
| aid | 账号ID |

**请求体**：
```json
{ "flags": { "key": "value" } }
```

**响应 data**：账号对象

---

## 平台 /platform

### 获取平台列表

`GET /platform/list`

**查询参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| enable | bool | 筛选已启用 |
| article | bool | 筛选支持文章 |
| graph_text | bool | 筛选支持图文 |
| video | bool | 筛选支持视频 |

**响应 data**：平台数组

### 切换平台状态

`POST /platform/state/:pid`

| 路径参数 | 说明 |
|---------|------|
| pid | 平台ID |

**响应 data**：平台对象

---

## 平台配置 /platSet

### 创建配置

`POST /platSet/create`

**请求体**：
```json
{
  "name": "配置名称",
  "description": "描述",
  "platform_id": 1,
  "setting": {}
}
```

**响应 data**：配置对象

### 更新配置

`POST /platSet/update/:sid`

| 路径参数 | 说明 |
|---------|------|
| sid | 配置ID |

**请求体**：同创建

**响应 data**：配置对象

### 设为默认配置

`POST /platSet/default/:sid`

| 路径参数 | 说明 |
|---------|------|
| sid | 配置ID |

**响应 data**：配置对象

### 启用/禁用配置

`POST /platSet/enable/:sid`

| 路径参数 | 说明 |
|---------|------|
| sid | 配置ID |

**响应 data**：配置对象

### 删除配置

`DELETE /platSet/delete/:sid`

| 路径参数 | 说明 |
|---------|------|
| sid | 配置ID |

**响应 data**：配置对象

### 获取平台配置列表

`GET /platSet/list/:pid`

| 路径参数 | 说明 |
|---------|------|
| pid | 平台ID |

**响应 data**：配置数组

### setting 参数说明

`setting` 是一个 JSON 对象，必须包含 `platType` 字段标识平台类型，其余字段根据平台不同而不同。

部分平台根据发布类型（文章/图文/视频）使用不同的 Setting 结构，具体见各平台说明。

**通用子对象 TimerPublish（定时发布）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| enable | bool | 是否启用定时发布 |
| timer | string | 定时时间，格式 `"2025-04-25 15:54:00"` |

---

#### wechat 微信公众号

**文章/图文使用 Setting：**

| 字段 | 类型 | 说明 |
|------|------|------|
| author | string | 作者 |
| link | string | 原文链接 |
| leave | bool | 开启留言，默认 true |
| origin | bool | 声明原创，默认 false |
| reprint | bool | 快捷转载，origin 为 true 时可设置 |
| publishType | string | 发表类型：`"mass"` 群发 / `"publish"` 发布 |
| collection | string | 合集（文章/图文/视频合集不能重复） |
| source | uint | 创作来源：0 不声明，1 内容由AI生成，2 素材来源官方媒体/网络新闻，3 内容剧情演绎仅供娱乐，4 个人观点仅供参考，5 健康医疗分享仅供参考，6 投资观点仅供参考，7 无需声明 |
| timerPublish | TimerPublish | 定时发布（可选时间：当前+5分钟 ~ 7天） |

**视频使用 VideoSetting（继承 Setting 全部字段，额外增加）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| materTitle | string | 素材标题 |
| barrage | bool | 弹幕 |
| barrageCheck | uint | 弹幕权限：0 所有用户，1 已关注用户，2 已关注7天及以上用户 |
| turn2Channel | bool | 发表后转为视频号视频 |
| adTrans | uint | 广告过渡：0 不设置，1~6 不同广告过渡语 |

---

#### wechat-video 微信视频号

| 字段 | 类型 | 说明 |
|------|------|------|
| location | string | 位置，默认 `"auto"` |
| collection | string | 合集 |
| linkType | uint | 链接类型：0 不设置，1 公众号文章，2 红包封面 |
| linkAddr | string | 链接地址 |
| music | string | 音乐 |
| activity | string | 活动 |
| origin | bool | 声明原创（仅视频） |
| timerPublish | TimerPublish | 定时发布 |

---

#### toutiaohao 今日头条

**文章使用 Setting：**

| 字段 | 类型 | 说明 |
|------|------|------|
| location | string | 位置 |
| placeAD | bool | 投放广告 |
| starter | bool | 头条首发 |
| collection | string | 合集（设置了合集不能定时发布） |
| syncPublish | bool | 同时发布微头条 |
| source | uint | 创作声明：0 不声明，1 取材网络，3 个人观点仅供参考，4 引用AI，5 虚构演绎故事经历，6 投资观点仅供参考，7 健康医疗分享仅供参考 |
| timerPublish | TimerPublish | 定时发布（可选时间：当前+2小时 ~ 7天） |

**图文使用 GTSetting（继承 Setting 全部字段，额外增加）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| openBgm | bool | 开启配乐 |

**视频使用 VideoSetting（独立结构）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| gtEnable | bool | 视频生成图文 |
| gtSyncPub | bool | 生成图文与视频同时发布 |
| collection | string | 合集 |
| stickers | array | 互动贴纸 |
| source | uint | 创作声明：0 不声明，1 取自站外，3 自行拍摄，4 AI生成，5 虚构演绎故事经历，6 投资观点仅供参考，7 健康医疗分享仅供参考 |
| link | string | 扩展链接 |
| lookScope | uint | 谁可以看：0 公开，1 粉丝，2 自己 |
| timerPublish | TimerPublish | 定时发布 |

---

#### douyin 抖音

| 字段 | 类型 | 说明 |
|------|------|------|
| activity | string | 添加活动奖励 |
| music | string | 音乐 |
| label | string | 标签（位置：带货模式/打卡模式，团购：全国，影视演绎，小程序） |
| location | string | 位置/商品/影视演艺/小程序/标记万物 |
| hotspot | string | 关联热点 |
| collection | string | 合集 |
| allowSave | bool | 允许他人保存，默认 true |
| lookScope | uint | 谁可以看：0 公开，1 好友，2 自己 |
| timerPublish | TimerPublish | 定时发布 |

---

#### kuaishou 快手

| 字段 | 类型 | 说明 |
|------|------|------|
| music | string | 添加音乐（仅图文） |
| linkApplet | string | 小程序链接 |
| source | uint | 作品声明：0 不声明，1 内容为AI生成，2 演绎情节仅供娱乐，3 个人观点仅供参考，4 素材来源于网络 |
| collection | string | 合集 |
| location | string | 位置 |
| sameFrame | bool | 允许别人跟我拍同框，默认 true |
| download | bool | 允许下载此作品，默认 true |
| sameCity | bool | 作品展示在同城页，默认 true |
| lookScope | uint | 谁可以看：0 公开，1 好友，2 自己 |
| timerPublish | TimerPublish | 定时发布 |

---

#### xiaohongshu 小红书

| 字段 | 类型 | 说明 |
|------|------|------|
| location | string | 位置 |
| collection | string | 合集 |
| group | string | 群聊 |
| mark | object/null | 标记：`{"user": true, "search": "搜索内容"}`，user 为 true 标记用户，false 标记地点 |
| origin | bool | 声明原创，默认 false |
| source | uint | 作品声明：0 不声明，1 虚构演绎仅供娱乐，2 笔记含AI合成内容，3 已在正文中自主标注，4 自主拍摄，5 来源转载 |
| reprint | string | 来源转载的来源媒体（source 为 5 时填写） |
| lookScope | uint | 谁可以看：0 公开，1 好友，2 自己 |
| timerPublish | TimerPublish | 定时发布 |

---

#### weishi 微视

| 字段 | 类型 | 说明 |
|------|------|------|
| source | uint | 作品声明：0 不声明，1 该内容由AI生成，2 剧情演绎仅供娱乐，3 个人观点仅供参考，4 取材网络谨慎甄别 |
| lookScope | uint | 谁可以看：0 公开，1 自己 |
| timerPublish | TimerPublish | 定时发布 |

---

#### bilibili 哔哩哔哩

**视频/图文使用 Setting：**

| 字段 | 类型 | 说明 |
|------|------|------|
| reprint | string | 转载来源，为空表示自制 |
| partition | string | 分区 |
| creation | bool | 是否允许二创 |
| public | bool | 是否公开可见 |
| source | uint | 创作声明：1 使用AI合成技术，2 含有危险行为，3 仅供娱乐，4 可能引人不适，5 理性适度消费，6 个人观点仅供参考 |
| dynamic | string | 粉丝动态（支持 @提及） |
| timerPublish | TimerPublish | 定时发布 |

**文章使用 ArticleSetting：**

| 字段 | 类型 | 说明 |
|------|------|------|
| classify | string | 专栏分类 |
| origin | bool | 声明原创，默认 false |
| headerImg | string | 头图 |
| labels | string | 标签，最多10个 |
| collection | string | 合集 |
| public | bool | 是否公开可见 |
| timerPublish | TimerPublish | 定时发布 |

---

#### omtencent 企鹅号

| 字段 | 类型 | 说明 |
|------|------|------|
| classify | string | 分类 |
| labels | string | 标签（多个用 `/` 分割） |
| activity | string | 活动 |
| source | uint | 自主声明 |
| timerPublish | TimerPublish | 定时发布 |

---

#### acfun A站

| 字段 | 类型 | 说明 |
|------|------|------|
| classify | string | 分区（格式：`"一级分区/二级分区"`） |
| labels | string | 标签（最多5个） |
| origin | bool | 类型：true 原创 / false 转载 |
| reprint | string | 转载来源（原创或文章不需要） |
| dynamic | string | 粉丝动态 |
| timerPublish | TimerPublish | 定时发布（可选时间：当前+4小时 ~ 14天） |

---

#### baijiahao 百家号

| 字段 | 类型 | 说明 |
|------|------|------|
| watermark | uint8 | 水印（仅视频）：0 不添加，1 添加水印，2 添加贴片 |
| location | string | 位置 |
| classify | string | 分类（格式：`"一级分类/二级分类"` 或 `"一级/二级/三级"`） |
| activity | string | 活动 |
| byAI | bool | AI创作声明 |
| timerPublish | TimerPublish | 定时发布（可选时间：当前+1小时 ~ 7天） |

---

#### zhihu 知乎

**文章/图文使用 Setting：**

| 字段 | 类型 | 说明 |
|------|------|------|
| question | string | 投稿至问题 |
| source | uint | 创作声明：0 无声明，1 包含剧透，2 包含医疗建议，3 虚构创作，4 包含理财内容，5 包含AI辅助创作 |
| topic | string | 文章话题，最多3个，多个用 `/` 分割 |
| collection | string | 专栏，为空表示不发布到专栏 |
| origin | uint | 内容来源：0 不设置，1 官方网站，2 新闻报道，3 电视媒体，4 纸质媒体 |

**视频使用 VideoSetting（继承 Setting 全部字段，额外增加）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| classify | string | 领域分类 |
| reprint | bool | true 转载 / false 原创 |
| timerPublish | TimerPublish | 定时发布 |

---

#### jianshuhao 简书

| 字段 | 类型 | 说明 |
|------|------|------|
| collection | string | 文集 |
| vetoReprint | bool | 禁止转载 |

---

#### juejin 掘金

| 字段 | 类型 | 说明 |
|------|------|------|
| classify | string | 分类 |
| tag | string | 标签（必填） |
| collection | string | 专栏 |
| topic | string | 话题 |
| group | string | 沸点圈子 |
| link | string | 沸点链接 |

---

#### sina 新浪微博

**视频/图文使用 Setting（继承 ArticleSetting 全部字段，额外增加）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| type | uint | 类型：0 原创，1 二创，2 转载 |
| classify | string | 分类（格式：`"栏目/分类"`） |
| stress | bool | 允许画重点，默认 true |
| location | string | 位置 |
| wait | int | 等待 X 秒后发布 |

**文章使用 ArticleSetting：**

| 字段 | 类型 | 说明 |
|------|------|------|
| collection | string | 专栏 |
| onlyFans | bool | 仅粉丝阅读全文，默认 true |
| lookScope | uint | 谁可以看：0 公开，1 粉丝（2 好友圈，3 自己 仅视频） |
| source | uint | 内容声明：0 不声明，1 内容由AI生成，2 内容为虚构演绎 |
| dynamic | string | 粉丝动态 |
| timerPublish | TimerPublish | 定时发布 |

---

#### csdn CSDN

**文章使用 Setting：**

| 字段 | 类型 | 说明 |
|------|------|------|
| labels | string | 标签，多个用 `/` 分割，最多7个 |
| collection | string | 分类专栏，多个用 `/` 分割，最多3个 |
| artType | uint | 文章类型：0 原创，1 转载，2 翻译 |
| originLink | string | 原文链接（转载必须，翻译可选） |
| backupGitCode | bool | 备份到 GitCode |
| lookScope | uint | 可见范围：0 全部，1 仅自己，2 粉丝可见，3 VIP可见 |
| activity | string | 参与活动 |
| topic | string | 话题 |
| timerPublish | TimerPublish | 定时发布（可选时间：当前+4小时 ~ 7天） |

**视频使用 VideoSetting：**

| 字段 | 类型 | 说明 |
|------|------|------|
| labels | string | 标签，多个用 `/` 分割，最多3个 |
| recommend | bool | 是否推荐 |

---

#### x X(Twitter)

| 字段 | 类型 | 说明 |
|------|------|------|
| consumerKey | string | API Consumer Key |
| consumerSecret | string | API Consumer Secret |
| replySettings | string | 回复权限：`"following"` / `"mentionedUsers"` / `"subscribers"` / `"verified"` |

---

#### tiktok TikTok

| 字段 | 类型 | 说明 |
|------|------|------|
| location | string | 位置 |
| lookScope | uint | 谁可以看：0 所有人，1 好友，2 自己 |
| comment | bool | 允许评论 |
| creation | bool | 二次创作内容 |
| reveal | bool | 披露作品内容 |
| yourBrand | bool | 你的品牌 |
| brandContent | bool | 品牌内容 |
| aigc | bool | AI生成的内容 |
| timerPublish | TimerPublish | 定时发布（可选时间：当前+2小时 ~ 30天） |

---

#### youtube YouTube

| 字段 | 类型 | 说明 |
|------|------|------|
| tags | string | 标签/关键词 |
| categoryId | string | 视频类别ID（如 `"22"` 人物与博客，`"28"` 科学与技术） |
| defaultLanguage | string | 默认语言（ISO 639-1，如 `"en"`, `"zh-CN"`） |
| localizations | string | 本地化信息（JSON 字符串） |
| embeddable | bool | 是否允许外部网站嵌入 |
| license | string | 许可证类型：`"youtube"` / `"creativeCommon"` |
| privacyStatus | string | 隐私状态：`"public"` / `"unlisted"` / `"private"` |
| publicStatsViewable | bool | 是否公开视频统计 |
| selfDeclaredMadeForKids | bool | 标记为面向儿童（COPPA） |
| containsSyntheticMedia | bool | 是否包含合成/虚拟内容 |
| recordingDate | string | 拍摄/录制日期（ISO 8601） |
| timerPublish | TimerPublish | 定时发布 |

---

#### pinduoduo 拼多多

| 字段 | 类型 | 说明 |
|------|------|------|
| goodsId | string | 商品ID |
| source | uint | 作品声明：0 不声明，1 内容由AI生成，2 内容取材网络，3 可能引人不适，4 虚构演绎仅供娱乐，5 危险行为请勿模仿 |
| timerPublish | TimerPublish | 定时发布（可选时间：当前+4小时 ~ 7天） |

---

## 内容 /article

### 获取内容列表

`GET /article/list`

**查询参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| status | uint8 | 状态筛选 |
| current | int | 页码 |
| size | int | 每页条数 |

**响应 data**：`{ "list": [], "pager": {} }`

### 获取内容详情

`GET /article/get/:rid`

| 路径参数 | 说明 |
|---------|------|
| rid | 内容ID |

**响应 data**：内容对象

### 创建文章

`POST /article/create`

**请求体**：PublishData 对象

**响应 data**：文章对象

### 创建图文

`POST /article/graphText`

**请求体**：PublishData 对象（无 thumb 时自动生成缩略图）

**响应 data**：图文对象

### 创建视频

`POST /article/video`

**请求体**：PublishData 对象（无 thumb 时自动生成缩略图）

**响应 data**：视频对象

### 更新内容

`POST /article/update/:rid`

| 路径参数 | 说明 |
|---------|------|
| rid | 内容ID |

**请求体**：PublishData 对象

**响应 data**：内容对象

### 删除内容

`DELETE /article/delete/:rid`

| 路径参数 | 说明 |
|---------|------|
| rid | 内容ID |

**响应 data**：内容对象

---

## SSE 发布 /sse

### 监听事件

`GET /sse/listen`

SSE 长连接，推送以下事件：

| 事件 | 说明 |
|------|------|
| heartbeat | 心跳（每分钟） |
| exit | 登录过期，需退出 |
| logout | 账号登出通知 |
| vip | VIP 到期提醒 |
| sync | 数据同步通知 |

### 发布文章

`POST /sse/article/:rid`

| 路径参数 | 说明 |
|---------|------|
| rid | 文章ID |

**请求体**：
```json
{
  "headless": false,
  "syncDraft": true,
  "postAccounts": [
    {
      "id": 1,
      "platName": "平台名",
      "settings": {}
    }
  ]
}
```

SSE 事件流：

| 事件 | 说明 |
|------|------|
| info | 发布进度信息 |
| success | 单个账号发布成功 |
| error | 发布错误 |
| finish | 发布任务完成，data: `{"msg":"","res":[true,false]}` |
| wait | 已有发布任务执行中 |
| vip | 非会员限制提示 |

### 发布图文

`POST /sse/graphText/:tid`

参数和事件同发布文章。

| 路径参数 | 说明 |
|---------|------|
| tid | 图文ID |

### 发布视频

`POST /sse/video/:vid`

参数和事件同发布文章。

| 路径参数 | 说明 |
|---------|------|
| vid | 视频ID |

---

## 发布记录 /record

### 获取发布记录列表

`GET /record/list`

**查询参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| status | int | 发布状态 |
| type | int | 发布类型 |
| current | int | 页码 |
| size | int | 每页条数 |

**响应 data**：`{ "list": [], "pager": {} }`

### 获取发布数据

`GET /record/publishData/:prid`

| 路径参数 | 说明 |
|---------|------|
| prid | 记录ID |

**响应 data**：PublishData 对象

### 获取发布详情

`GET /record/info/:prid`

| 路径参数 | 说明 |
|---------|------|
| prid | 记录ID |

**响应 data**：发布详情数组

### 删除发布记录

`DELETE /record/delete/:prid`

| 路径参数 | 说明 |
|---------|------|
| prid | 记录ID |

支持批量删除：`DELETE /record/delete/:prid?prids=1,2,3`

**响应 data**：记录对象 或 `null`（批量）

---

## 仪表盘 /home

### 仪表盘数据

`GET /home/dashboard`

**响应 data**：
```json
{
  "article_data": {},
  "account_data": {},
  "publish_data": {},
  "record_data": {}
}
```

### 概览数据

`GET /home/overview`

**响应 data**：
```json
{
  "record_data": {},
  "chart_data": []
}
```

### 内容排行

`GET /home/content`

**查询参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| limit | int | 返回条数 |
| order | int | 排序方式 |

**响应 data**：内容排行数组

---

## 图片 /pix

### 上传图片

`POST /pix/upload`

**请求体**：`multipart/form-data`，字段 `file`

**响应 data**：图片访问 URL 字符串

### 获取图片

`GET /pix/:filename`

无需认证头，直接返回图片二进制内容。

### 获取图片 Base64

`GET /pix/:filename/raw`

**响应 data**：`data:{mime};base64,{data}` 格式字符串

---

## 文档转换 /pandoc

### 导入文档

`POST /pandoc/input`

VIP功能。将文档文件转换为 Markdown。

**请求体**：
```json
{ "input": "文件路径" }
```

支持格式：`.md` `.rst` `.mediawiki` `.wiki` `.docx` `.html` `.htm`

**响应 data**：Markdown 字符串

### 导出文档

`POST /pandoc/output`

VIP功能。将 Markdown 导出为指定格式文件。

**请求体**：
```json
{
  "markdown": "Markdown 内容",
  "output": "输出文件路径（含扩展名）"
}
```

**响应 data**：`null`

---

## YouTube 授权 /youtube

### 获取授权URL

`POST /youtube/authUrl`

**请求体**：
```json
{ "proxy": "代理地址（可选）" }
```

**响应 data**：OAuth 授权 URL 字符串

### 授权回调

`GET /youtubeAuthCallback`

无需认证。OAuth 回调接口，返回 HTML 页面。

### 授权结果

`GET /oauth/authRes`

**响应 data**：
```json
{
  "authing": false,
  "error": ""
}
```

---

## X (Twitter) 授权 /x

### 获取授权URL

`POST /x/authUrl`

**请求体**：
```json
{ "proxy": "代理地址（可选）" }
```

**响应 data**：OAuth 授权 URL 字符串，或配置说明链接

---

## 代理管理 /proxy

### 创建代理

`POST /proxy/create`

VIP功能。

**请求体**：
```json
{
  "scheme": "http",
  "host": "127.0.0.1",
  "port": 8080,
  "username": "",
  "password": ""
}
```

**响应 data**：代理对象

### 获取代理列表

`GET /proxy/list`

**查询参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |
| health | bool | 筛选健康代理 |

**响应 data**：代理数组

### 更新代理

`POST /proxy/update/:pid`

VIP功能。

| 路径参数 | 说明 |
|---------|------|
| pid | 代理ID |

**请求体**：
```json
{
  "scheme": "http",
  "host": "127.0.0.1",
  "port": 8080,
  "username": "",
  "password": ""
}
```

**响应 data**：代理对象

### 删除代理

`DELETE /proxy/delete/:pid`

| 路径参数 | 说明 |
|---------|------|
| pid | 代理ID |

**响应 data**：代理对象

### 检测代理

`POST /proxy/check/:pid`

VIP功能。检测单个代理延迟。

| 路径参数 | 说明 |
|---------|------|
| pid | 代理ID |

**响应 data**：代理对象（含 delay）

### 批量检测代理

`POST /proxy/batch/check`

VIP功能。检测所有代理健康状态。

**响应 data**：
```json
{
  "total": 10,
  "health": 8,
  "timeout": 1,
  "failed": 1
}
```

### 导入代理

`POST /proxy/import`

VIP功能。从文本文件导入代理列表。

**请求体**：
```json
{ "path": "文件路径" }
```

**响应 data**：导入的代理数组

### 导出代理

`GET /proxy/export`

VIP功能。

**查询参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| path | string | 导出文件路径 |

**响应 data**：文件路径字符串

### 可用代理

`GET /proxy/unused/:pid`

VIP功能。

| 路径参数 | 说明 |
|---------|------|
| pid | 平台ID |

**响应 data**：代理数组

---

## 数据导入导出 /sys

### 导出数据

`POST /sys/export`

导出账号、内容、平台配置、代理数据。

**请求体**：
```json
{ "path": "导出文件路径" }
```

**响应 data**：导出结果

### 导入数据

`POST /sys/import`

**请求体**：
```json
{ "path": "导入文件路径" }
```

**响应 data**：导入结果

---

## 支付 /pay

### 创建订单

`POST /pay/tradeCreate`

**请求体**：
```json
{ "paymentType": "支付类型" }
```

**响应 data**：订单信息

### 查询订单

`POST /pay/tradeQuery`

**请求体**：
```json
{ "outTradeNo": "订单号" }
```

**响应 data**：支付结果字符串 或 `null`

### 关闭订单

`POST /pay/tradeClose`

**响应 data**：`null`
