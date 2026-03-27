package main

import "sort"

// FieldDef 描述一个平台 setting 字段
type FieldDef struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // "string", "bool", "uint", "int", "object", "array"
	Required    bool        `json:"required"`
	Description string      `json:"description"`
	Default     any         `json:"default,omitempty"`
	Options     []OptionDef `json:"options,omitempty"`
}

// OptionDef 描述枚举选项
type OptionDef struct {
	Value       any    `json:"value"`
	Description string `json:"description"`
}

// 内容类型常量
const (
	ContentAll       = "*"
	ContentArticle   = "article"
	ContentGraphText = "graph_text"
	ContentVideo     = "video"
)

// getSchema 查询指定平台和内容类型的 setting schema
// 优先精确匹配 contentType，未命中则 fallback 到 "*"
func getSchema(platType, contentType string) ([]FieldDef, bool) {
	ct, ok := platformSchemas[platType]
	if !ok {
		return nil, false
	}
	if fields, ok := ct[contentType]; ok {
		return fields, true
	}
	if fields, ok := ct[ContentAll]; ok {
		return fields, true
	}
	return nil, false
}

// getRequiredFields 返回指定平台和内容类型的所有必填字段名
func getRequiredFields(platType, contentType string) []string {
	fields, ok := getSchema(platType, contentType)
	if !ok {
		return nil
	}
	var required []string
	for _, f := range fields {
		if f.Required {
			required = append(required, f.Name)
		}
	}
	return required
}

// getSupportedPlatTypes 返回所有支持的 platType 列表（排序）
func getSupportedPlatTypes() []string {
	keys := make([]string, 0, len(platformSchemas))
	for k := range platformSchemas {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// --- 通用字段片段，供多个平台复用 ---

var timerPublishField = FieldDef{
	Name:        "timerPublish",
	Type:        "object",
	Description: "定时发布配置（settings 中的字段）。对象格式：{\"enable\": true, \"timer\": \"2025-04-25 15:54:00\"}。enable=true 启用定时发布，timer 为发布时间（格式：YYYY-MM-DD HH:mm:ss）",
}

var sourceFieldFactory = func(desc string, options []OptionDef) FieldDef {
	return FieldDef{
		Name:        "source",
		Type:        "uint",
		Description: desc,
		Default:     0,
		Options:     options,
	}
}

// lookScopePublic 公开/好友/自己 三档可见范围（抖音、快手等）
var lookScopePublic = FieldDef{
	Name:        "lookScope",
	Type:        "uint",
	Description: "谁可以看",
	Default:     0,
	Options: []OptionDef{
		{Value: 0, Description: "公开"},
		{Value: 1, Description: "好友"},
		{Value: 2, Description: "自己"},
	},
}

// --- 平台 Schema 定义 ---

var platformSchemas = map[string]map[string][]FieldDef{

	// ==================== 微信公众号 ====================
	"wechat": {
		ContentArticle:   wechatBaseSetting(),
		ContentGraphText: wechatBaseSetting(),
		ContentVideo: append(wechatBaseSetting(),
			FieldDef{Name: "materTitle", Type: "string", Description: "素材标题，不填写则默认为视频素材文件名"},
			FieldDef{Name: "turn2Channel", Type: "bool", Description: "转为视频号视频（请确保公众号已绑定视频号）", Default: false},
			FieldDef{Name: "barrage", Type: "bool", Description: "弹幕", Default: true},
			FieldDef{Name: "barrageCheck", Type: "uint", Description: "谁可以发弹幕", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "所有用户"},
				{Value: 1, Description: "已关注用户"},
				{Value: 2, Description: "已关注7天及以上用户"},
			}},
		),
	},

	// ==================== 微信视频号 ====================
	"wechat-video": {
		ContentAll: {
			{Name: "location", Type: "string", Description: "地理位置，不输入表示不显示位置，输入 auto 表示由平台自动设置，其他内容表示自定义位置", Default: "auto"},
			{Name: "collection", Type: "string", Description: "合集（请确保合集已在平台创建）"},
			{Name: "linkType", Type: "uint", Description: "链接类型", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "不设置"},
				{Value: 1, Description: "公众号文章"},
				{Value: 2, Description: "红包封面"},
			}},
			{Name: "linkAddr", Type: "string", Description: "链接地址（仅在设置链接类型后生效）"},
			{Name: "activity", Type: "string", Description: "参与活动"},
			{Name: "origin", Type: "bool", Description: "声明原创（仅视频，请确保账号有声明原创权限）", Default: false},
			{Name: "music", Type: "string", Description: "音乐（仅对图文类型生效）"},
			timerPublishField,
		},
	},

	// ==================== 今日头条 ====================
	"toutiaohao": {
		ContentArticle:   toutiaohaoArticleSetting(),
		ContentGraphText: toutiaohaoGtSetting(),
		ContentVideo:     toutiaohaoVideoSetting(),
	},

	// ==================== 抖音 ====================
	"douyin": {
		ContentAll: {
			{Name: "activity", Type: "string", Description: "平台活动（输入活动名称，多个活动用 / 分割）"},
			{Name: "collection", Type: "string", Description: "合集（请确保合集已在平台创建并通过审核）"},
			{Name: "label", Type: "string", Description: "扩展信息标签", Options: []OptionDef{
				{Value: "位置", Description: "位置"},
				{Value: "影视演艺", Description: "影视演艺"},
				{Value: "小程序", Description: "小程序"},
				{Value: "标记万物", Description: "标记万物"},
			}},
			{Name: "location", Type: "string", Description: "标签值：位置标签值不输入将默认选择猜你想加或第一个选项；团购标签值必须按「范围/商品名/推广标题」格式输入"},
			{Name: "hotspot", Type: "string", Description: "关联热点（将选择搜索到的第一个结果）"},
			{Name: "allowSave", Type: "bool", Description: "允许他人保存", Default: true},
			{Name: "lookScope", Type: "uint", Description: "可见范围（加入合集后该项配置不生效）", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "公开"},
				{Value: 1, Description: "好友"},
				{Value: 2, Description: "自己"},
			}},
			{Name: "music", Type: "string", Description: "音乐（仅对图文类型生效）"},
			timerPublishField,
		},
	},

	// ==================== 快手 ====================
	"kuaishou": {
		ContentAll: {
			{Name: "linkApplet", Type: "string", Description: "关联小程序（前往快手 APP 小程序复制页面链接）"},
			sourceFieldFactory("作品声明", []OptionDef{
				{Value: 0, Description: "不声明"},
				{Value: 1, Description: "内容为AI生成"},
				{Value: 2, Description: "演绎情节仅供娱乐"},
				{Value: 3, Description: "个人观点仅供参考"},
				{Value: 4, Description: "素材来源于网络"},
			}),
			{Name: "collection", Type: "string", Description: "合集（请确保合集已在平台创建）"},
			{Name: "location", Type: "string", Description: "位置（将选择搜索到的第一个结果）"},
			{Name: "sameFrame", Type: "bool", Description: "允许别人跟我拍同框", Default: true},
			{Name: "download", Type: "bool", Description: "允许下载此作品", Default: true},
			{Name: "sameCity", Type: "bool", Description: "作品展示在同城页", Default: true},
			lookScopePublic,
			{Name: "music", Type: "string", Description: "音乐（仅对图文类型生效）"},
			timerPublishField,
		},
	},

	// ==================== 小红书 ====================
	"xiaohongshu": {
		ContentAll: {
			{Name: "location", Type: "string", Description: "位置（将选择搜索到的第一个结果）"},
			{Name: "collection", Type: "string", Description: "合集（请确保合集已在平台创建）"},
			{Name: "group", Type: "string", Description: "关联群聊（请确保群聊已在平台创建）"},
			{Name: "mark", Type: "object", Description: "标记：{\"user\": true, \"search\": \"搜索内容\"}，user 为 true 标记用户，false 标记地点"},
			{Name: "origin", Type: "bool", Description: "声明原创", Default: false},
			sourceFieldFactory("创作来源", []OptionDef{
				{Value: 0, Description: "不声明"},
				{Value: 1, Description: "虚构演绎仅供娱乐"},
				{Value: 2, Description: "笔记含AI合成内容"},
				{Value: 3, Description: "已在正文中自主标注"},
				{Value: 4, Description: "自主拍摄"},
				{Value: 5, Description: "来源转载"},
			}),
			{Name: "reprint", Type: "string", Description: "关联媒体（仅在创作来源为来源转载时生效）"},
			{Name: "lookScope", Type: "uint", Description: "可见范围", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "公开"},
				{Value: 1, Description: "互关好友"},
				{Value: 2, Description: "自己"},
			}},
			timerPublishField,
		},
	},

	// ==================== 微视 ====================
	"weishi": {
		ContentAll: {
			sourceFieldFactory("作品声明", []OptionDef{
				{Value: 0, Description: "不声明"},
				{Value: 1, Description: "该内容由AI生成"},
				{Value: 2, Description: "剧情演绎仅供娱乐"},
				{Value: 3, Description: "个人观点仅供参考"},
				{Value: 4, Description: "取材网络谨慎甄别"},
			}),
			{Name: "lookScope", Type: "uint", Description: "谁可以看", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "公开"},
				{Value: 1, Description: "自己"},
			}},
			timerPublishField,
		},
	},

	// ==================== 哔哩哔哩 ====================
	"bilibili": {
		ContentVideo:     bilibiliVideoSetting(),
		ContentGraphText: bilibiliVideoSetting(),
		ContentArticle: {
			{Name: "classify", Type: "string", Description: "专栏分类，格式：\"一级/二级\"，如 \"生活/日常\"", Default: "生活/日常"},
			{Name: "origin", Type: "bool", Description: "声明原创", Default: false},
			{Name: "headerImg", Type: "string", Description: "头图（本地图片路径）"},
			{Name: "labels", Type: "string", Description: "文章标签，多个用 / 分割，最多10个"},
			{Name: "collection", Type: "string", Description: "合集（请确保合集已在平台创建）"},
			{Name: "public", Type: "bool", Description: "是否公开可见，false 则仅自己可见", Default: true},
			timerPublishField,
		},
	},

	// ==================== 企鹅号 ====================
	"omtencent": {
		ContentArticle:   omtencentArticleSetting(),
		ContentGraphText: omtencentCommonSetting(),
		ContentVideo:     omtencentVideoSetting(),
	},

	// ==================== A站 ====================
	"acfun": {
		ContentAll: {
			{Name: "origin", Type: "bool", Description: "原创声明：true 原创 / false 转载", Default: true},
			{Name: "classify", Type: "string", Required: true, Description: "分区（格式：\"一级分区/二级分区\"），文章和视频分区不同，必填"},
			{Name: "labels", Type: "string", Description: "文章标签，多个用 / 分割，最多6个（仅对文章生效，视频标签从描述话题自动提取）"},
			{Name: "reprint", Type: "string", Description: "转载来源（origin 为 false 且发布视频时必填）"},
			{Name: "dynamic", Type: "string", Description: "粉丝动态"},
			timerPublishField,
		},
	},

	// ==================== 百家号 ====================
	"baijiahao": {
		ContentArticle:   baijiahaoBaseSetting(),
		ContentGraphText: baijiahaoBaseSetting(),
		ContentVideo: append(baijiahaoBaseSetting(),
			FieldDef{Name: "watermark", Type: "uint", Description: "视频水印（仅对视频生效）", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "不添加水印"},
				{Value: 1, Description: "添加水印"},
				{Value: 2, Description: "添加贴片"},
			}},
		),
	},

	// ==================== 知乎 ====================
	"zhihu": {
		ContentArticle:   zhihuArticleSetting(),
		ContentGraphText: zhihuArticleSetting(),
		ContentVideo:     zhihuVideoSetting(),
	},

	// ==================== 简书 ====================
	"jianshuhao": {
		ContentAll: {
			{Name: "collection", Type: "string", Description: "文集（请确保文集已在平台创建）"},
			{Name: "vetoReprint", Type: "bool", Description: "禁止转载", Default: false},
		},
	},

	// ==================== 掘金 ====================
	"juejin": {
		ContentArticle: {
			{Name: "classify", Type: "string", Description: "分类", Options: []OptionDef{
				{Value: "后端", Description: "后端"},
				{Value: "前端", Description: "前端"},
				{Value: "Android", Description: "Android"},
				{Value: "IOS", Description: "IOS"},
				{Value: "人工智能", Description: "人工智能"},
				{Value: "开发工具", Description: "开发工具"},
				{Value: "代码人生", Description: "代码人生"},
				{Value: "阅读", Description: "阅读"},
			}},
			{Name: "tag", Type: "string", Required: true, Description: "标签（必填，当前仅支持单个标签）"},
			{Name: "collection", Type: "string", Description: "专栏（请确保专栏已在平台创建并通过审核）"},
			{Name: "topic", Type: "string", Description: "话题（选择与文章相关的话题，否则可能被平台移除）"},
		},
		ContentGraphText: {
			{Name: "group", Type: "string", Description: "沸点圈子（仅在发布图文/沸点时生效）"},
			{Name: "link", Type: "string", Description: "沸点链接（仅在发布图文/沸点时生效）"},
		},
	},

	// ==================== 新浪微博 ====================
	"sina": {
		ContentArticle:   sinaArticleSetting(),
		ContentGraphText: sinaGtSetting(),
		ContentVideo:     sinaVideoSetting(),
	},

	// ==================== CSDN ====================
	"csdn": {
		ContentArticle: {
			{Name: "labels", Type: "string", Required: true, Description: "标签，多个用 / 分割，最多7个"},
			{Name: "collection", Type: "string", Description: "分类专栏，多个用 / 分割，最多3个（请确保专栏已在平台创建）"},
			{Name: "artType", Type: "uint", Description: "文章类型", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "原创"},
				{Value: 1, Description: "转载"},
				{Value: 2, Description: "翻译"},
			}},
			{Name: "originLink", Type: "string", Description: "原文链接（转载必须设置，翻译可选；转载不能设置可见范围为VIP可见）"},
			{Name: "backupGitCode", Type: "bool", Description: "同时备份到 GitCode", Default: false},
			{Name: "lookScope", Type: "uint", Description: "可见范围", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "全部"},
				{Value: 1, Description: "仅自己"},
				{Value: 2, Description: "粉丝可见"},
				{Value: 3, Description: "VIP可见"},
			}},
			{Name: "activity", Type: "string", Description: "参与活动"},
			{Name: "topic", Type: "string", Description: "话题"},
			timerPublishField,
		},
		ContentVideo: {
			{Name: "labels", Type: "string", Required: true, Description: "标签，多个用 / 分割，最多3个"},
			{Name: "recommend", Type: "bool", Description: "是否推荐", Default: false},
		},
	},

	// ==================== X (Twitter) ====================
	"x": {
		ContentAll: {
			{Name: "consumerKey", Type: "string", Description: "Consumer API Key"},
			{Name: "consumerSecret", Type: "string", Description: "Consumer API Secret"},
			{Name: "replySettings", Type: "string", Description: "回复权限", Options: []OptionDef{
				{Value: "following", Description: "关注的人"},
				{Value: "mentionedUsers", Description: "提及的用户"},
				{Value: "subscribers", Description: "订阅者"},
				{Value: "verified", Description: "已认证用户"},
			}},
		},
	},

	// ==================== TikTok ====================
	"tiktok": {
		ContentVideo: {
			{Name: "location", Type: "string", Description: "位置（将选择搜索到的第一个结果）"},
			{Name: "lookScope", Type: "uint", Description: "谁可以看（品牌内容作品不支持设置为自己）", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "所有人"},
				{Value: 1, Description: "好友"},
				{Value: 2, Description: "自己"},
			}},
			{Name: "comment", Type: "bool", Description: "允许评论", Default: true},
			{Name: "creation", Type: "bool", Description: "二次创作内容（超过60秒的视频无法进行合拍/抢镜及拼接）", Default: true},
			{Name: "reveal", Type: "bool", Description: "披露作品内容（告知此作品推广的是品牌、商品或服务）", Default: false},
			{Name: "yourBrand", Type: "bool", Description: "你的品牌（你正在推广自己或自己的业务）", Default: false},
			{Name: "brandContent", Type: "bool", Description: "品牌内容（与某品牌是品牌合作关系）", Default: false},
			{Name: "aigc", Type: "bool", Description: "AI生成的内容", Default: false},
			timerPublishField,
		},
	},

	// ==================== YouTube ====================
	"youtube": {
		ContentVideo: {
			{Name: "tags", Type: "string", Description: "标签/关键词，多个用 / 分割（可帮助提高搜索发现率）"},
			{Name: "categoryId", Type: "string", Description: "视频类别", Default: "22", Options: []OptionDef{
				{Value: "1", Description: "电影与动画"},
				{Value: "2", Description: "汽车与交通工具"},
				{Value: "10", Description: "音乐"},
				{Value: "15", Description: "宠物与动物"},
				{Value: "17", Description: "运动"},
				{Value: "18", Description: "短片"},
				{Value: "19", Description: "旅行与活动"},
				{Value: "20", Description: "游戏"},
				{Value: "21", Description: "视频博客"},
				{Value: "22", Description: "人物与博客"},
				{Value: "23", Description: "喜剧"},
				{Value: "24", Description: "娱乐"},
				{Value: "25", Description: "新闻与政治"},
				{Value: "26", Description: "教育"},
				{Value: "27", Description: "科学与技术"},
				{Value: "28", Description: "非营利组织与社会公益"},
				{Value: "29", Description: "电视剧"},
				{Value: "30", Description: "预告片"},
			}},
			{Name: "defaultLanguage", Type: "string", Description: "默认语言（ISO 639-1 代码，如 \"zh\", \"en\"）", Default: "zh"},
			{Name: "embeddable", Type: "bool", Description: "是否允许外部网站嵌入播放", Default: true},
			{Name: "license", Type: "string", Description: "许可协议", Default: "youtube", Options: []OptionDef{
				{Value: "youtube", Description: "标准 YouTube 许可协议"},
				{Value: "creativeCommon", Description: "Creative Commons (CC BY)"},
			}},
			{Name: "privacyStatus", Type: "string", Description: "隐私状态", Default: "private", Options: []OptionDef{
				{Value: "public", Description: "公开"},
				{Value: "unlisted", Description: "不公开（仅持链接可见）"},
				{Value: "private", Description: "私有"},
			}},
			{Name: "publicStatsViewable", Type: "bool", Description: "是否向公众显示观看数、点赞数等统计数据", Default: true},
			{Name: "selfDeclaredMadeForKids", Type: "bool", Description: "是否标记为\"面向儿童\"（COPPA 要求）", Default: false},
			{Name: "containsSyntheticMedia", Type: "bool", Description: "是否包含合成/虚拟人物或 AI 生成内容", Default: false},
			{Name: "recordingDate", Type: "string", Description: "视频拍摄/录制日期（ISO 8601 格式，如 \"2025-01-01T00:00:00Z\"）"},
			{Name: "localizations", Type: "string", Description: "多语言信息（JSON 格式），如 {\"en\":{\"title\":\"English title\",\"description\":\"English desc\"}}"},
			timerPublishField,
		},
	},

	// ==================== 拼多多 ====================
	"pinduoduo": {
		ContentVideo: {
			{Name: "goodsId", Type: "string", Description: "推广商品的ID（示例：12345678）"},
			sourceFieldFactory("内容声明", []OptionDef{
				{Value: 0, Description: "不声明"},
				{Value: 1, Description: "内容由AI生成"},
				{Value: 2, Description: "内容取材网络"},
				{Value: 3, Description: "可能引人不适"},
				{Value: 4, Description: "虚构演绎仅供娱乐"},
				{Value: 5, Description: "危险行为请勿模仿"},
			}),
			timerPublishField,
		},
	},
}

// --- 辅助函数：生成共享的基础字段切片 ---

// wechatBaseSetting 微信公众号通用设置（文章/图文/视频公共部分）
func wechatBaseSetting() []FieldDef {
	return []FieldDef{
		{Name: "author", Type: "string", Description: "作者名称"},
		{Name: "link", Type: "string", Description: "原文链接（点击阅读原文的跳转链接）"},
		{Name: "leave", Type: "bool", Description: "开启留言", Default: true},
		{Name: "origin", Type: "bool", Description: "声明原创", Default: true},
		{Name: "reprint", Type: "bool", Description: "快捷转载（origin 为 true 时可设置）", Default: true},
		{Name: "publishType", Type: "string", Description: "发表类型", Default: "mass", Options: []OptionDef{
			{Value: "mass", Description: "群发"},
			{Value: "publish", Description: "发布"},
		}},
		{Name: "collection", Type: "string", Description: "合集（文章/图文/视频合集不能重复，请确保合集已在平台创建）"},
		sourceFieldFactory("创作来源", []OptionDef{
			{Value: 0, Description: "不声明"},
			{Value: 1, Description: "内容由AI生成"},
			{Value: 2, Description: "素材来源官方媒体/网络新闻"},
			{Value: 3, Description: "内容剧情演绎仅供娱乐"},
			{Value: 4, Description: "个人观点仅供参考"},
			{Value: 5, Description: "健康医疗分享仅供参考"},
			{Value: 6, Description: "投资观点仅供参考"},
			{Value: 7, Description: "无需声明"},
		}),
		timerPublishField,
	}
}

// toutiaohaoArticleSetting 头条号文章设置
func toutiaohaoArticleSetting() []FieldDef {
	return []FieldDef{
		{Name: "location", Type: "string", Description: "位置（平台仅支持地级市，如：广州）"},
		{Name: "collection", Type: "string", Description: "合集（设置了合集不能定时发布，将直接发布）"},
		{Name: "starter", Type: "bool", Description: "头条首发（72小时内仅在头条发布内容）", Default: false},
		{Name: "placeAD", Type: "bool", Description: "投放广告", Default: true},
		{Name: "syncPublish", Type: "bool", Description: "同时发布微头条", Default: true},
		sourceFieldFactory("作品声明", []OptionDef{
			{Value: 0, Description: "不声明"},
			{Value: 1, Description: "取材网络"},
			{Value: 3, Description: "个人观点仅供参考"},
			{Value: 4, Description: "引用AI"},
			{Value: 5, Description: "虚构演绎故事经历"},
			{Value: 6, Description: "投资观点仅供参考"},
			{Value: 7, Description: "健康医疗分享仅供参考"},
		}),
		timerPublishField,
	}
}

// toutiaohaoGtSetting 头条号图文设置（含 openBgm）
func toutiaohaoGtSetting() []FieldDef {
	return []FieldDef{
		{Name: "location", Type: "string", Description: "位置（将选择搜索到的第一个结果）"},
		{Name: "starter", Type: "bool", Description: "头条首发（72小时内仅在头条发布内容）", Default: false},
		sourceFieldFactory("作品声明", []OptionDef{
			{Value: 0, Description: "不声明"},
			{Value: 1, Description: "取材网络"},
			{Value: 3, Description: "个人观点仅供参考"},
			{Value: 4, Description: "引用AI"},
			{Value: 5, Description: "虚构演绎故事经历"},
		}),
		{Name: "openBgm", Type: "bool", Description: "开启配乐，开启后可在小视频场景分发", Default: false},
		timerPublishField,
	}
}

// toutiaohaoVideoSetting 头条号视频设置
func toutiaohaoVideoSetting() []FieldDef {
	return []FieldDef{
		{Name: "gtEnable", Type: "bool", Description: "视频生成图文", Default: true},
		{Name: "gtSyncPub", Type: "bool", Description: "生成图文与视频同时发布（否则仅保存草稿）", Default: true},
		{Name: "collection", Type: "string", Description: "合集（竖版视频不支持，设置了合集不能定时发布）"},
		sourceFieldFactory("作品声明", []OptionDef{
			{Value: 0, Description: "不声明"},
			{Value: 1, Description: "取自站外"},
			{Value: 3, Description: "自行拍摄"},
			{Value: 4, Description: "AI生成"},
			{Value: 5, Description: "虚构演绎故事经历"},
			{Value: 6, Description: "投资观点仅供参考"},
			{Value: 7, Description: "健康医疗分享仅供参考"},
		}),
		{Name: "link", Type: "string", Description: "扩展链接（竖版视频不支持）"},
		{Name: "lookScope", Type: "uint", Description: "可见范围（竖版视频不支持）", Default: 0, Options: []OptionDef{
			{Value: 0, Description: "公开"},
			{Value: 1, Description: "粉丝"},
			{Value: 2, Description: "自己"},
		}},
		timerPublishField,
	}
}

// bilibiliVideoSetting bilibili 视频/图文设置
func bilibiliVideoSetting() []FieldDef {
	return []FieldDef{
		{Name: "reprint", Type: "string", Description: "转载来源，为空表示自制，非空表示转载"},
		{Name: "partition", Type: "string", Description: "分区", Options: []OptionDef{
			{Value: "影视", Description: "影视"},
			{Value: "娱乐", Description: "娱乐"},
			{Value: "音乐", Description: "音乐"},
			{Value: "舞蹈", Description: "舞蹈"},
			{Value: "动画", Description: "动画"},
			{Value: "绘画", Description: "绘画"},
			{Value: "鬼畜", Description: "鬼畜"},
			{Value: "游戏", Description: "游戏"},
			{Value: "资讯", Description: "资讯"},
			{Value: "知识", Description: "知识"},
			{Value: "人工智能", Description: "人工智能"},
			{Value: "科技数码", Description: "科技数码"},
			{Value: "汽车", Description: "汽车"},
			{Value: "时尚美妆", Description: "时尚美妆"},
			{Value: "家装房产", Description: "家装房产"},
			{Value: "户外潮流", Description: "户外潮流"},
			{Value: "健身", Description: "健身"},
			{Value: "体育运动", Description: "体育运动"},
			{Value: "手工", Description: "手工"},
			{Value: "美食", Description: "美食"},
			{Value: "小剧场", Description: "小剧场"},
			{Value: "旅游出行", Description: "旅游出行"},
			{Value: "三农", Description: "三农"},
			{Value: "动物", Description: "动物"},
			{Value: "亲子", Description: "亲子"},
			{Value: "健康", Description: "健康"},
			{Value: "情感", Description: "情感"},
			{Value: "vlog", Description: "vlog"},
			{Value: "生活兴趣", Description: "生活兴趣"},
			{Value: "生活经验", Description: "生活经验"},
		}},
		{Name: "creation", Type: "bool", Description: "是否允许二创", Default: true},
		{Name: "public", Type: "bool", Description: "是否公开可见，false 则仅自己可见", Default: true},
		sourceFieldFactory("创作声明", []OptionDef{
			{Value: 0, Description: "不声明"},
			{Value: 1, Description: "该视频使用人工智能合成技术"},
			{Value: 2, Description: "视频内含有危险行为请勿轻易模仿"},
			{Value: 3, Description: "该内容仅供娱乐请勿过分解读"},
			{Value: 4, Description: "该内容可能引人不适请谨慎选择观看"},
			{Value: 5, Description: "请理性适度消费"},
			{Value: 6, Description: "个人观点仅供参考"},
		}),
		{Name: "dynamic", Type: "string", Description: "粉丝动态（支持 @提及）"},
		timerPublishField,
	}
}

// omtencentCommonSetting 企鹅号通用字段（用于图文）
func omtencentCommonSetting() []FieldDef {
	return []FieldDef{
		{Name: "labels", Type: "string", Description: "标签，多个用 / 分割，最多9个，每个标签最多8个字"},
		{Name: "activity", Type: "string", Description: "参与活动"},
		sourceFieldFactory("声明", []OptionDef{
			{Value: 0, Description: "不声明"},
			{Value: 1, Description: "该内容由AI生成"},
			{Value: 2, Description: "演绎情节仅供娱乐"},
			{Value: 3, Description: "取材网络谨慎甄别"},
			{Value: 4, Description: "个人观点仅供参考"},
			{Value: 5, Description: "旧闻"},
		}),
		timerPublishField,
	}
}

// omtencentArticleSetting 企鹅号文章设置（通用字段 + 文章分类）
func omtencentArticleSetting() []FieldDef {
	fields := omtencentCommonSetting()
	fields = append(fields, FieldDef{
		Name:        "classify",
		Type:        "string",
		Description: "文章分类",
		Options: []OptionDef{
			{Value: "财经", Description: "财经"}, {Value: "彩票", Description: "彩票"},
			{Value: "CBA", Description: "CBA"}, {Value: "宠物", Description: "宠物"},
			{Value: "创意", Description: "创意"}, {Value: "传媒", Description: "传媒"},
			{Value: "出国", Description: "出国"}, {Value: "电影", Description: "电影"},
			{Value: "动漫", Description: "动漫"}, {Value: "法律", Description: "法律"},
			{Value: "房产", Description: "房产"}, {Value: "佛教", Description: "佛教"},
			{Value: "搞笑", Description: "搞笑"}, {Value: "GIF", Description: "GIF"},
			{Value: "股票", Description: "股票"}, {Value: "环球", Description: "环球"},
			{Value: "互联网", Description: "互联网"}, {Value: "婚庆", Description: "婚庆"},
			{Value: "户外运动", Description: "户外运动"}, {Value: "家居", Description: "家居"},
			{Value: "减肥", Description: "减肥"}, {Value: "健康", Description: "健康"},
			{Value: "健身", Description: "健身"}, {Value: "教育", Description: "教育"},
			{Value: "基督教", Description: "基督教"}, {Value: "鸡汤", Description: "鸡汤"},
			{Value: "军事", Description: "军事"}, {Value: "科技", Description: "科技"},
			{Value: "科学", Description: "科学"}, {Value: "篮球", Description: "篮球"},
			{Value: "历史", Description: "历史"}, {Value: "旅游", Description: "旅游"},
			{Value: "美女", Description: "美女"}, {Value: "美食", Description: "美食"},
			{Value: "美图", Description: "美图"}, {Value: "命理", Description: "命理"},
			{Value: "民俗", Description: "民俗"}, {Value: "NBA", Description: "NBA"},
			{Value: "跑步", Description: "跑步"}, {Value: "汽车", Description: "汽车"},
			{Value: "情感", Description: "情感"}, {Value: "人物", Description: "人物"},
			{Value: "三农", Description: "三农"}, {Value: "社会", Description: "社会"},
			{Value: "生活百科", Description: "生活百科"}, {Value: "摄影", Description: "摄影"},
			{Value: "时尚", Description: "时尚"}, {Value: "时政", Description: "时政"},
			{Value: "数码", Description: "数码"}, {Value: "天气", Description: "天气"},
			{Value: "体育", Description: "体育"}, {Value: "文化", Description: "文化"},
			{Value: "文学", Description: "文学"}, {Value: "星座", Description: "星座"},
			{Value: "新闻", Description: "新闻"}, {Value: "意甲", Description: "意甲"},
			{Value: "英超", Description: "英超"}, {Value: "音乐", Description: "音乐"},
			{Value: "艺术", Description: "艺术"}, {Value: "游戏", Description: "游戏"},
			{Value: "育儿", Description: "育儿"}, {Value: "娱乐", Description: "娱乐"},
			{Value: "招聘", Description: "招聘"}, {Value: "职场", Description: "职场"},
			{Value: "中超", Description: "中超"}, {Value: "宗教", Description: "宗教"},
			{Value: "足球", Description: "足球"},
		},
	})
	return fields
}

// omtencentVideoSetting 企鹅号视频设置（通用字段 + 视频分类）
func omtencentVideoSetting() []FieldDef {
	fields := omtencentCommonSetting()
	fields = append(fields, FieldDef{
		Name:        "classify",
		Type:        "string",
		Description: "视频分类",
		Options: []OptionDef{
			{Value: "连载动画", Description: "连载动画"}, {Value: "原创动画", Description: "原创动画"},
			{Value: "动漫周边", Description: "动漫周边"}, {Value: "宅文化", Description: "宅文化"},
			{Value: "美食纪录片", Description: "美食纪录片"}, {Value: "自然纪录片", Description: "自然纪录片"},
			{Value: "历史纪录片", Description: "历史纪录片"}, {Value: "社会纪录片", Description: "社会纪录片"},
			{Value: "旅游纪录片", Description: "旅游纪录片"}, {Value: "军事纪录片", Description: "军事纪录片"},
			{Value: "人文纪录片", Description: "人文纪录片"}, {Value: "评测", Description: "评测"},
			{Value: "产品资讯", Description: "产品资讯"}, {Value: "智能生活", Description: "智能生活"},
			{Value: "数码其他", Description: "数码其他"}, {Value: "电影周边", Description: "电影周边"},
			{Value: "电影资讯", Description: "电影资讯"}, {Value: "电影剪辑", Description: "电影剪辑"},
			{Value: "微电影", Description: "微电影"}, {Value: "影评", Description: "影评"},
			{Value: "农业生产养殖", Description: "农业生产养殖"}, {Value: "农家生活", Description: "农家生活"},
			{Value: "农村政策", Description: "农村政策"}, {Value: "美味食谱", Description: "美味食谱"},
			{Value: "烘焙", Description: "烘焙"}, {Value: "美食猎奇", Description: "美食猎奇"},
			{Value: "吃播大胃王", Description: "吃播大胃王"}, {Value: "明星资讯", Description: "明星资讯"},
			{Value: "饭拍饭制", Description: "饭拍饭制"}, {Value: "篮球", Description: "篮球"},
			{Value: "足球", Description: "足球"}, {Value: "体育趣闻", Description: "体育趣闻"},
			{Value: "体育教学", Description: "体育教学"}, {Value: "综合体育", Description: "综合体育"},
			{Value: "功夫搏击", Description: "功夫搏击"}, {Value: "极限运动", Description: "极限运动"},
			{Value: "广场舞教学", Description: "广场舞教学"}, {Value: "广场舞欣赏", Description: "广场舞欣赏"},
			{Value: "新车速递", Description: "新车速递"}, {Value: "试驾评测", Description: "试驾评测"},
			{Value: "汽车资讯", Description: "汽车资讯"}, {Value: "用车", Description: "用车"},
			{Value: "玩车", Description: "玩车"}, {Value: "车模", Description: "车模"},
			{Value: "摩托车", Description: "摩托车"}, {Value: "汽车其他", Description: "汽车其他"},
			{Value: "爆笑动物", Description: "爆笑动物"}, {Value: "爆笑恶搞", Description: "爆笑恶搞"},
			{Value: "相声小品", Description: "相声小品"}, {Value: "鬼畜", Description: "鬼畜"},
			{Value: "奇趣", Description: "奇趣"}, {Value: "段子剧", Description: "段子剧"},
			{Value: "神回复", Description: "神回复"}, {Value: "熊孩子", Description: "熊孩子"},
			{Value: "糗事", Description: "糗事"}, {Value: "爆笑原创", Description: "爆笑原创"},
			{Value: "广告", Description: "广告"}, {Value: "创意", Description: "创意"},
			{Value: "自拍", Description: "自拍"}, {Value: "公益短片", Description: "公益短片"},
			{Value: "青春", Description: "青春"}, {Value: "绝活", Description: "绝活"},
			{Value: "法制", Description: "法制"}, {Value: "交通", Description: "交通"},
			{Value: "正能量", Description: "正能量"}, {Value: "国内时政", Description: "国内时政"},
			{Value: "国际时政", Description: "国际时政"}, {Value: "国际社会", Description: "国际社会"},
			{Value: "社会百态", Description: "社会百态"}, {Value: "人物访谈", Description: "人物访谈"},
			{Value: "奇闻趣事", Description: "奇闻趣事"}, {Value: "孕产", Description: "孕产"},
			{Value: "喂养", Description: "喂养"}, {Value: "早教", Description: "早教"},
			{Value: "萌宝", Description: "萌宝"}, {Value: "保健护理", Description: "保健护理"},
			{Value: "科学实验", Description: "科学实验"}, {Value: "太空探索", Description: "太空探索"},
			{Value: "科普", Description: "科普"}, {Value: "科学其他", Description: "科学其他"},
			{Value: "舞台剧", Description: "舞台剧"}, {Value: "魔术", Description: "魔术"},
			{Value: "戏曲", Description: "戏曲"}, {Value: "其他综艺", Description: "其他综艺"},
			{Value: "音乐真人秀", Description: "音乐真人秀"}, {Value: "喜剧节目", Description: "喜剧节目"},
			{Value: "脱口秀", Description: "脱口秀"}, {Value: "真人秀", Description: "真人秀"},
			{Value: "情感节目", Description: "情感节目"}, {Value: "儿歌", Description: "儿歌"},
			{Value: "益智", Description: "益智"}, {Value: "少儿动画", Description: "少儿动画"},
			{Value: "少儿节目", Description: "少儿节目"}, {Value: "读书", Description: "读书"},
			{Value: "艺术", Description: "艺术"}, {Value: "宗教", Description: "宗教"},
			{Value: "历史", Description: "历史"}, {Value: "文化", Description: "文化"},
			{Value: "电视剧剪辑", Description: "电视剧剪辑"}, {Value: "电视剧片花", Description: "电视剧片花"},
			{Value: "剧集周边", Description: "剧集周边"}, {Value: "网络剧", Description: "网络剧"},
			{Value: "健康", Description: "健康"}, {Value: "休闲", Description: "休闲"},
			{Value: "健身", Description: "健身"}, {Value: "家居", Description: "家居"},
			{Value: "生活窍门", Description: "生活窍门"}, {Value: "风水命理", Description: "风水命理"},
			{Value: "心灵", Description: "心灵"}, {Value: "彩票", Description: "彩票"},
			{Value: "宠物", Description: "宠物"}, {Value: "摄影", Description: "摄影"},
			{Value: "两性", Description: "两性"}, {Value: "星座", Description: "星座"},
			{Value: "情感", Description: "情感"}, {Value: "消费", Description: "消费"},
			{Value: "投资", Description: "投资"}, {Value: "理财", Description: "理财"},
			{Value: "产经", Description: "产经"}, {Value: "股市", Description: "股市"},
			{Value: "金融", Description: "金融"}, {Value: "房地产", Description: "房地产"},
			{Value: "财经人物", Description: "财经人物"}, {Value: "收藏", Description: "收藏"},
			{Value: "公司", Description: "公司"}, {Value: "创业", Description: "创业"},
			{Value: "经济", Description: "经济"}, {Value: "考试", Description: "考试"},
			{Value: "技能教学", Description: "技能教学"}, {Value: "校园内外", Description: "校园内外"},
			{Value: "人生课堂", Description: "人生课堂"}, {Value: "语言学习", Description: "语言学习"},
			{Value: "知识百科", Description: "知识百科"}, {Value: "演讲", Description: "演讲"},
			{Value: "公开课", Description: "公开课"}, {Value: "职场教育", Description: "职场教育"},
			{Value: "武器装备", Description: "武器装备"}, {Value: "战争历史", Description: "战争历史"},
			{Value: "军情解读", Description: "军情解读"}, {Value: "军事节目", Description: "军事节目"},
			{Value: "美妆", Description: "美妆"}, {Value: "潮流奢品", Description: "潮流奢品"},
			{Value: "男士时尚", Description: "男士时尚"}, {Value: "T台秀场", Description: "T台秀场"},
			{Value: "穿搭", Description: "穿搭"}, {Value: "时尚资讯", Description: "时尚资讯"},
			{Value: "时尚大片", Description: "时尚大片"}, {Value: "舞蹈达人", Description: "舞蹈达人"},
			{Value: "舞蹈教学", Description: "舞蹈教学"}, {Value: "舞蹈工作室", Description: "舞蹈工作室"},
			{Value: "明星热舞", Description: "明星热舞"}, {Value: "互联网", Description: "互联网"},
			{Value: "科技前沿", Description: "科技前沿"}, {Value: "科技奇趣", Description: "科技奇趣"},
			{Value: "航空航天", Description: "航空航天"}, {Value: "机械", Description: "机械"},
			{Value: "演唱会", Description: "演唱会"}, {Value: "音乐节目", Description: "音乐节目"},
			{Value: "恶搞音乐", Description: "恶搞音乐"}, {Value: "翻唱", Description: "翻唱"},
			{Value: "演奏", Description: "演奏"}, {Value: "影视音乐", Description: "影视音乐"},
			{Value: "明星MV", Description: "明星MV"}, {Value: "达人MV", Description: "达人MV"},
			{Value: "喊麦", Description: "喊麦"}, {Value: "音乐牛人", Description: "音乐牛人"},
			{Value: "手机游戏", Description: "手机游戏"}, {Value: "网络游戏", Description: "网络游戏"},
			{Value: "单机游戏", Description: "单机游戏"}, {Value: "游戏节目", Description: "游戏节目"},
			{Value: "达人解说", Description: "达人解说"}, {Value: "电竞赛事", Description: "电竞赛事"},
			{Value: "游戏展会", Description: "游戏展会"}, {Value: "游戏周边", Description: "游戏周边"},
			{Value: "旅行攻略", Description: "旅行攻略"}, {Value: "旅行趣闻", Description: "旅行趣闻"},
			{Value: "旅途风光", Description: "旅途风光"},
		},
	})
	return fields
}

// zhihuArticleSetting 知乎文章/图文设置
func zhihuArticleSetting() []FieldDef {
	return []FieldDef{
		{Name: "origin", Type: "uint", Description: "内容来源", Default: 0, Options: []OptionDef{
			{Value: 0, Description: "不设置"},
			{Value: 1, Description: "官方网站"},
			{Value: 2, Description: "新闻报道"},
			{Value: 3, Description: "电视媒体"},
			{Value: 4, Description: "纸质媒体"},
		}},
		{Name: "question", Type: "string", Description: "投稿至问题（将选择搜索到的第一个结果）"},
		sourceFieldFactory("创作声明", []OptionDef{
			{Value: 0, Description: "无声明"},
			{Value: 1, Description: "包含剧透"},
			{Value: 2, Description: "包含医疗建议"},
			{Value: 3, Description: "虚构创作"},
			{Value: 4, Description: "包含理财内容"},
			{Value: 5, Description: "包含AI辅助创作"},
		}),
		{Name: "topic", Type: "string", Description: "文章话题，最多3个，多个用 / 分割"},
		{Name: "collection", Type: "string", Description: "专栏（为空表示不发布到专栏）"},
	}
}

// zhihuVideoSetting 知乎视频设置
func zhihuVideoSetting() []FieldDef {
	return []FieldDef{
		{Name: "origin", Type: "uint", Description: "内容来源", Default: 0, Options: []OptionDef{
			{Value: 0, Description: "不设置"},
			{Value: 1, Description: "官方网站"},
			{Value: 2, Description: "新闻报道"},
			{Value: 3, Description: "电视媒体"},
			{Value: 4, Description: "纸质媒体"},
		}},
		{Name: "classify", Type: "string", Required: true, Description: "所属领域（发布视频必须设置）", Options: []OptionDef{
			{Value: "人文", Description: "人文"}, {Value: "体育竞技", Description: "体育竞技"},
			{Value: "健康医学", Description: "健康医学"}, {Value: "其他", Description: "其他"},
			{Value: "军事", Description: "军事"}, {Value: "动漫", Description: "动漫"},
			{Value: "娱乐", Description: "娱乐"}, {Value: "宠物", Description: "宠物"},
			{Value: "家居生活", Description: "家居生活"}, {Value: "家用电器", Description: "家用电器"},
			{Value: "影视", Description: "影视"}, {Value: "心理学", Description: "心理学"},
			{Value: "情感", Description: "情感"}, {Value: "故事", Description: "故事"},
			{Value: "教育", Description: "教育"}, {Value: "数码", Description: "数码"},
			{Value: "旅行", Description: "旅行"}, {Value: "时尚穿搭", Description: "时尚穿搭"},
			{Value: "母婴亲子", Description: "母婴亲子"}, {Value: "汽车", Description: "汽车"},
			{Value: "法律", Description: "法律"}, {Value: "游戏电竞", Description: "游戏电竞"},
			{Value: "社会/时政", Description: "社会/时政"}, {Value: "社会学", Description: "社会学"},
			{Value: "科学工程", Description: "科学工程"}, {Value: "科技互联网", Description: "科技互联网"},
			{Value: "经济与管理", Description: "经济与管理"}, {Value: "美妆个护", Description: "美妆个护"},
			{Value: "美食", Description: "美食"}, {Value: "职场", Description: "职场"},
			{Value: "艺术", Description: "艺术"}, {Value: "运动健身", Description: "运动健身"},
			{Value: "音乐", Description: "音乐"},
		}},
		{Name: "reprint", Type: "bool", Description: "视频类型：true 转载 / false 原创", Default: false},
	}
}

// sinaArticleSetting 新浪微博文章设置
func sinaArticleSetting() []FieldDef {
	return []FieldDef{
		{Name: "collection", Type: "string", Description: "专栏（请确保专栏已在平台创建）"},
		{Name: "onlyFans", Type: "bool", Description: "仅粉丝阅读全文", Default: true},
		{Name: "lookScope", Type: "uint", Description: "可见范围", Default: 0, Options: []OptionDef{
			{Value: 0, Description: "公开"},
			{Value: 1, Description: "粉丝"},
		}},
		sourceFieldFactory("内容声明", []OptionDef{
			{Value: 0, Description: "不声明"},
			{Value: 1, Description: "内容由AI生成"},
			{Value: 2, Description: "内容为虚构演绎"},
		}),
		{Name: "dynamic", Type: "string", Description: "微博动态（支持话题/@用户/微博表情，如：#话题名称# @用户名称 内容[666]）"},
		timerPublishField,
	}
}

// sinaGtSetting 新浪微博图文（微博）设置
func sinaGtSetting() []FieldDef {
	return []FieldDef{
		{Name: "location", Type: "string", Description: "地点（将选择搜索到的第一个结果）"},
		{Name: "lookScope", Type: "uint", Description: "可见范围", Default: 0, Options: []OptionDef{
			{Value: 0, Description: "公开"},
			{Value: 1, Description: "粉丝"},
			{Value: 2, Description: "好友圈"},
			{Value: 3, Description: "自己"},
		}},
		sourceFieldFactory("内容声明", []OptionDef{
			{Value: 0, Description: "不声明"},
			{Value: 1, Description: "内容由AI生成"},
			{Value: 2, Description: "内容为虚构演绎"},
		}),
		timerPublishField,
	}
}

// sinaVideoSetting 新浪微博视频设置
func sinaVideoSetting() []FieldDef {
	return []FieldDef{
		{Name: "wait", Type: "int", Description: "等待 X 秒后发布（用于处理频繁发布触发的人机验证，在等待时间内完成验证）"},
		{Name: "type", Type: "uint", Description: "视频类型", Default: 0, Options: []OptionDef{
			{Value: 0, Description: "原创"},
			{Value: 1, Description: "二创"},
			{Value: 2, Description: "转载"},
		}},
		{Name: "classify", Type: "string", Description: "视频分类（格式：\"栏目/分类\"）"},
		{Name: "stress", Type: "bool", Description: "允许画重点", Default: true},
		{Name: "collection", Type: "string", Description: "合集（允许加入多个合集，多个用 / 分割）"},
		{Name: "location", Type: "string", Description: "地点（将选择搜索到的第一个结果）"},
		{Name: "lookScope", Type: "uint", Description: "可见范围", Default: 0, Options: []OptionDef{
			{Value: 0, Description: "公开"},
			{Value: 1, Description: "粉丝"},
			{Value: 2, Description: "好友圈"},
			{Value: 3, Description: "自己"},
		}},
		sourceFieldFactory("内容声明", []OptionDef{
			{Value: 0, Description: "不声明"},
			{Value: 1, Description: "内容由AI生成"},
			{Value: 2, Description: "内容为虚构演绎"},
		}),
		timerPublishField,
	}
}

// baijiahaoBaseSetting 百家号通用设置
func baijiahaoBaseSetting() []FieldDef {
	return []FieldDef{
		{Name: "location", Type: "string", Description: "位置（将选择搜索到的第一个结果，仅对视频和图文生效）"},
		{Name: "classify", Type: "string", Description: "分类（格式：\"一级分类/二级分类\" 或 \"一级/二级/三级\"）"},
		{Name: "activity", Type: "string", Description: "活动（仅支持单个活动）"},
		{Name: "byAI", Type: "bool", Description: "AI创作声明", Default: false},
		timerPublishField,
	}
}
