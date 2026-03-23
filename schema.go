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
	Description: "定时发布，对象格式：{\"enable\": true, \"timer\": \"2025-04-25 15:54:00\"}",
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
			FieldDef{Name: "materTitle", Type: "string", Description: "素材标题"},
			FieldDef{Name: "barrage", Type: "bool", Description: "弹幕"},
			FieldDef{Name: "barrageCheck", Type: "uint", Description: "弹幕权限", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "所有用户"},
				{Value: 1, Description: "已关注用户"},
				{Value: 2, Description: "已关注7天及以上用户"},
			}},
			FieldDef{Name: "turn2Channel", Type: "bool", Description: "发表后转为视频号视频"},
			FieldDef{Name: "adTrans", Type: "uint", Description: "广告过渡", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "不设置"},
				{Value: 1, Description: "广告过渡语1"},
				{Value: 2, Description: "广告过渡语2"},
				{Value: 3, Description: "广告过渡语3"},
				{Value: 4, Description: "广告过渡语4"},
				{Value: 5, Description: "广告过渡语5"},
				{Value: 6, Description: "广告过渡语6"},
			}},
		),
	},

	// ==================== 微信视频号 ====================
	"wechat-video": {
		ContentAll: {
			{Name: "location", Type: "string", Description: "位置", Default: "auto"},
			{Name: "collection", Type: "string", Description: "合集"},
			{Name: "linkType", Type: "uint", Description: "链接类型", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "不设置"},
				{Value: 1, Description: "公众号文章"},
				{Value: 2, Description: "红包封面"},
			}},
			{Name: "linkAddr", Type: "string", Description: "链接地址"},
			{Name: "music", Type: "string", Description: "音乐"},
			{Name: "activity", Type: "string", Description: "活动"},
			{Name: "origin", Type: "bool", Description: "声明原创（仅视频）", Default: false},
			timerPublishField,
		},
	},

	// ==================== 今日头条 ====================
	"toutiaohao": {
		ContentArticle: toutiaohaoBaseSetting(),
		ContentGraphText: append(toutiaohaoBaseSetting(),
			FieldDef{Name: "openBgm", Type: "bool", Description: "开启配乐"},
		),
		ContentVideo: {
			{Name: "gtEnable", Type: "bool", Description: "视频生成图文"},
			{Name: "gtSyncPub", Type: "bool", Description: "生成图文与视频同时发布"},
			{Name: "collection", Type: "string", Description: "合集"},
			{Name: "stickers", Type: "array", Description: "互动贴纸"},
			sourceFieldFactory("创作声明", []OptionDef{
				{Value: 0, Description: "不声明"},
				{Value: 1, Description: "取自站外"},
				{Value: 3, Description: "自行拍摄"},
				{Value: 4, Description: "AI生成"},
				{Value: 5, Description: "虚构演绎故事经历"},
				{Value: 6, Description: "投资观点仅供参考"},
				{Value: 7, Description: "健康医疗分享仅供参考"},
			}),
			{Name: "link", Type: "string", Description: "扩展链接"},
			lookScopePublic,
			timerPublishField,
		},
	},

	// ==================== 抖音 ====================
	"douyin": {
		ContentAll: {
			{Name: "activity", Type: "string", Description: "添加活动奖励"},
			{Name: "music", Type: "string", Description: "音乐"},
			{Name: "label", Type: "string", Description: "标签（位置：带货模式/打卡模式，团购：全国，影视演绎，小程序）"},
			{Name: "location", Type: "string", Description: "位置/商品/影视演艺/小程序/标记万物"},
			{Name: "hotspot", Type: "string", Description: "关联热点"},
			{Name: "collection", Type: "string", Description: "合集"},
			{Name: "allowSave", Type: "bool", Description: "允许他人保存", Default: true},
			lookScopePublic,
			timerPublishField,
		},
	},

	// ==================== 快手 ====================
	"kuaishou": {
		ContentAll: {
			{Name: "music", Type: "string", Description: "添加音乐（仅图文）"},
			{Name: "linkApplet", Type: "string", Description: "小程序链接"},
			sourceFieldFactory("作品声明", []OptionDef{
				{Value: 0, Description: "不声明"},
				{Value: 1, Description: "内容为AI生成"},
				{Value: 2, Description: "演绎情节仅供娱乐"},
				{Value: 3, Description: "个人观点仅供参考"},
				{Value: 4, Description: "素材来源于网络"},
			}),
			{Name: "collection", Type: "string", Description: "合集"},
			{Name: "location", Type: "string", Description: "位置"},
			{Name: "sameFrame", Type: "bool", Description: "允许别人跟我拍同框", Default: true},
			{Name: "download", Type: "bool", Description: "允许下载此作品", Default: true},
			{Name: "sameCity", Type: "bool", Description: "作品展示在同城页", Default: true},
			lookScopePublic,
			timerPublishField,
		},
	},

	// ==================== 小红书 ====================
	"xiaohongshu": {
		ContentAll: {
			{Name: "location", Type: "string", Description: "位置"},
			{Name: "collection", Type: "string", Description: "合集"},
			{Name: "group", Type: "string", Description: "群聊"},
			{Name: "mark", Type: "object", Description: "标记：{\"user\": true, \"search\": \"搜索内容\"}，user 为 true 标记用户，false 标记地点"},
			{Name: "origin", Type: "bool", Description: "声明原创", Default: false},
			sourceFieldFactory("作品声明", []OptionDef{
				{Value: 0, Description: "不声明"},
				{Value: 1, Description: "虚构演绎仅供娱乐"},
				{Value: 2, Description: "笔记含AI合成内容"},
				{Value: 3, Description: "已在正文中自主标注"},
				{Value: 4, Description: "自主拍摄"},
				{Value: 5, Description: "来源转载"},
			}),
			{Name: "reprint", Type: "string", Description: "来源转载的来源媒体（source 为 5 时填写）"},
			lookScopePublic,
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
			{Name: "classify", Type: "string", Description: "专栏分类"},
			{Name: "origin", Type: "bool", Description: "声明原创", Default: false},
			{Name: "headerImg", Type: "string", Description: "头图"},
			{Name: "labels", Type: "string", Description: "标签，最多10个"},
			{Name: "collection", Type: "string", Description: "合集"},
			{Name: "public", Type: "bool", Description: "是否公开可见", Default: true},
			timerPublishField,
		},
	},

	// ==================== 企鹅号 ====================
	"omtencent": {
		ContentAll: {
			{Name: "classify", Type: "string", Description: "分类"},
			{Name: "labels", Type: "string", Description: "标签（多个用 / 分割）"},
			{Name: "activity", Type: "string", Description: "活动"},
			sourceFieldFactory("自主声明", nil),
			timerPublishField,
		},
	},

	// ==================== A站 ====================
	"acfun": {
		ContentAll: {
			{Name: "classify", Type: "string", Required: true, Description: "分区（格式：\"一级分区/二级分区\"）"},
			{Name: "labels", Type: "string", Description: "标签（最多5个）"},
			{Name: "origin", Type: "bool", Description: "类型：true 原创 / false 转载"},
			{Name: "reprint", Type: "string", Description: "转载来源（原创或文章不需要，转载时必填）"},
			{Name: "dynamic", Type: "string", Description: "粉丝动态"},
			timerPublishField,
		},
	},

	// ==================== 百家号 ====================
	"baijiahao": {
		ContentArticle:   baijiahaoBaseSetting(),
		ContentGraphText: baijiahaoBaseSetting(),
		ContentVideo: append(baijiahaoBaseSetting(),
			FieldDef{Name: "watermark", Type: "uint", Description: "水印", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "不添加"},
				{Value: 1, Description: "添加水印"},
				{Value: 2, Description: "添加贴片"},
			}},
		),
	},

	// ==================== 知乎 ====================
	"zhihu": {
		ContentArticle:   zhihuBaseSetting(),
		ContentGraphText: zhihuBaseSetting(),
		ContentVideo: append(zhihuBaseSetting(),
			FieldDef{Name: "classify", Type: "string", Description: "领域分类"},
			FieldDef{Name: "reprint", Type: "bool", Description: "true 转载 / false 原创"},
			timerPublishField,
		),
	},

	// ==================== 简书 ====================
	"jianshuhao": {
		ContentAll: {
			{Name: "collection", Type: "string", Description: "文集"},
			{Name: "vetoReprint", Type: "bool", Description: "禁止转载"},
		},
	},

	// ==================== 掘金 ====================
	"juejin": {
		ContentAll: {
			{Name: "classify", Type: "string", Description: "分类"},
			{Name: "tag", Type: "string", Required: true, Description: "标签（必填）"},
			{Name: "collection", Type: "string", Description: "专栏"},
			{Name: "topic", Type: "string", Description: "话题"},
			{Name: "group", Type: "string", Description: "沸点圈子"},
			{Name: "link", Type: "string", Description: "沸点链接"},
		},
	},

	// ==================== 新浪微博 ====================
	"sina": {
		ContentArticle: sinaArticleSetting(),
		ContentVideo: append(sinaArticleSetting(),
			FieldDef{Name: "type", Type: "uint", Description: "类型", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "原创"},
				{Value: 1, Description: "二创"},
				{Value: 2, Description: "转载"},
			}},
			FieldDef{Name: "classify", Type: "string", Description: "分类（格式：\"栏目/分类\"）"},
			FieldDef{Name: "stress", Type: "bool", Description: "允许画重点", Default: true},
			FieldDef{Name: "location", Type: "string", Description: "位置"},
			FieldDef{Name: "wait", Type: "int", Description: "等待 X 秒后发布"},
		),
		ContentGraphText: append(sinaArticleSetting(),
			FieldDef{Name: "type", Type: "uint", Description: "类型", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "原创"},
				{Value: 1, Description: "二创"},
				{Value: 2, Description: "转载"},
			}},
			FieldDef{Name: "classify", Type: "string", Description: "分类（格式：\"栏目/分类\"）"},
			FieldDef{Name: "stress", Type: "bool", Description: "允许画重点", Default: true},
			FieldDef{Name: "location", Type: "string", Description: "位置"},
			FieldDef{Name: "wait", Type: "int", Description: "等待 X 秒后发布"},
		),
	},

	// ==================== CSDN ====================
	"csdn": {
		ContentArticle: {
			{Name: "labels", Type: "string", Description: "标签，多个用 / 分割，最多7个"},
			{Name: "collection", Type: "string", Description: "分类专栏，多个用 / 分割，最多3个"},
			{Name: "artType", Type: "uint", Description: "文章类型", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "原创"},
				{Value: 1, Description: "转载"},
				{Value: 2, Description: "翻译"},
			}},
			{Name: "originLink", Type: "string", Description: "原文链接（转载必须，翻译可选）"},
			{Name: "backupGitCode", Type: "bool", Description: "备份到 GitCode"},
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
			{Name: "labels", Type: "string", Description: "标签，多个用 / 分割，最多3个"},
			{Name: "recommend", Type: "bool", Description: "是否推荐"},
		},
	},

	// ==================== X (Twitter) ====================
	"x": {
		ContentAll: {
			{Name: "consumerKey", Type: "string", Description: "API Consumer Key"},
			{Name: "consumerSecret", Type: "string", Description: "API Consumer Secret"},
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
		ContentAll: {
			{Name: "location", Type: "string", Description: "位置"},
			{Name: "lookScope", Type: "uint", Description: "谁可以看", Default: 0, Options: []OptionDef{
				{Value: 0, Description: "所有人"},
				{Value: 1, Description: "好友"},
				{Value: 2, Description: "自己"},
			}},
			{Name: "comment", Type: "bool", Description: "允许评论"},
			{Name: "creation", Type: "bool", Description: "二次创作内容"},
			{Name: "reveal", Type: "bool", Description: "披露作品内容"},
			{Name: "yourBrand", Type: "bool", Description: "你的品牌"},
			{Name: "brandContent", Type: "bool", Description: "品牌内容"},
			{Name: "aigc", Type: "bool", Description: "AI生成的内容"},
			timerPublishField,
		},
	},

	// ==================== YouTube ====================
	"youtube": {
		ContentAll: {
			{Name: "tags", Type: "string", Description: "标签/关键词"},
			{Name: "categoryId", Type: "string", Description: "视频类别ID（如 \"22\" 人物与博客，\"28\" 科学与技术）"},
			{Name: "defaultLanguage", Type: "string", Description: "默认语言（ISO 639-1，如 \"en\", \"zh-CN\"）"},
			{Name: "localizations", Type: "string", Description: "本地化信息（JSON 字符串）"},
			{Name: "embeddable", Type: "bool", Description: "是否允许外部网站嵌入"},
			{Name: "license", Type: "string", Description: "许可证类型", Options: []OptionDef{
				{Value: "youtube", Description: "YouTube 标准许可"},
				{Value: "creativeCommon", Description: "Creative Commons"},
			}},
			{Name: "privacyStatus", Type: "string", Description: "隐私状态", Default: "public", Options: []OptionDef{
				{Value: "public", Description: "公开"},
				{Value: "unlisted", Description: "不公开列出"},
				{Value: "private", Description: "私享"},
			}},
			{Name: "publicStatsViewable", Type: "bool", Description: "是否公开视频统计"},
			{Name: "selfDeclaredMadeForKids", Type: "bool", Description: "标记为面向儿童（COPPA）"},
			{Name: "containsSyntheticMedia", Type: "bool", Description: "是否包含合成/虚拟内容"},
			{Name: "recordingDate", Type: "string", Description: "拍摄/录制日期（ISO 8601）"},
			timerPublishField,
		},
	},

	// ==================== 拼多多 ====================
	"pinduoduo": {
		ContentAll: {
			{Name: "goodsId", Type: "string", Description: "商品ID"},
			sourceFieldFactory("作品声明", []OptionDef{
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

func wechatBaseSetting() []FieldDef {
	return []FieldDef{
		{Name: "author", Type: "string", Description: "作者"},
		{Name: "link", Type: "string", Description: "原文链接"},
		{Name: "leave", Type: "bool", Description: "开启留言", Default: true},
		{Name: "origin", Type: "bool", Description: "声明原创", Default: false},
		{Name: "reprint", Type: "bool", Description: "快捷转载，origin 为 true 时可设置"},
		{Name: "publishType", Type: "string", Description: "发表类型", Options: []OptionDef{
			{Value: "mass", Description: "群发"},
			{Value: "publish", Description: "发布"},
		}},
		{Name: "collection", Type: "string", Description: "合集（文章/图文/视频合集不能重复）"},
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

func toutiaohaoBaseSetting() []FieldDef {
	return []FieldDef{
		{Name: "location", Type: "string", Description: "位置"},
		{Name: "placeAD", Type: "bool", Description: "投放广告"},
		{Name: "starter", Type: "bool", Description: "头条首发"},
		{Name: "collection", Type: "string", Description: "合集（设置了合集不能定时发布）"},
		{Name: "syncPublish", Type: "bool", Description: "同时发布微头条"},
		sourceFieldFactory("创作声明", []OptionDef{
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

func bilibiliVideoSetting() []FieldDef {
	return []FieldDef{
		{Name: "reprint", Type: "string", Description: "转载来源，为空表示自制"},
		{Name: "partition", Type: "string", Description: "分区"},
		{Name: "creation", Type: "bool", Description: "是否允许二创"},
		{Name: "public", Type: "bool", Description: "是否公开可见", Default: true},
		sourceFieldFactory("创作声明", []OptionDef{
			{Value: 1, Description: "使用AI合成技术"},
			{Value: 2, Description: "含有危险行为"},
			{Value: 3, Description: "仅供娱乐"},
			{Value: 4, Description: "可能引人不适"},
			{Value: 5, Description: "理性适度消费"},
			{Value: 6, Description: "个人观点仅供参考"},
		}),
		{Name: "dynamic", Type: "string", Description: "粉丝动态（支持 @提及）"},
		timerPublishField,
	}
}

func zhihuBaseSetting() []FieldDef {
	return []FieldDef{
		{Name: "question", Type: "string", Description: "投稿至问题"},
		sourceFieldFactory("创作声明", []OptionDef{
			{Value: 0, Description: "无声明"},
			{Value: 1, Description: "包含剧透"},
			{Value: 2, Description: "包含医疗建议"},
			{Value: 3, Description: "虚构创作"},
			{Value: 4, Description: "包含理财内容"},
			{Value: 5, Description: "包含AI辅助创作"},
		}),
		{Name: "topic", Type: "string", Description: "文章话题，最多3个，多个用 / 分割"},
		{Name: "collection", Type: "string", Description: "专栏，为空表示不发布到专栏"},
		{Name: "origin", Type: "uint", Description: "内容来源", Default: 0, Options: []OptionDef{
			{Value: 0, Description: "不设置"},
			{Value: 1, Description: "官方网站"},
			{Value: 2, Description: "新闻报道"},
			{Value: 3, Description: "电视媒体"},
			{Value: 4, Description: "纸质媒体"},
		}},
	}
}

func sinaArticleSetting() []FieldDef {
	return []FieldDef{
		{Name: "collection", Type: "string", Description: "专栏"},
		{Name: "onlyFans", Type: "bool", Description: "仅粉丝阅读全文", Default: true},
		{Name: "lookScope", Type: "uint", Description: "谁可以看", Default: 0, Options: []OptionDef{
			{Value: 0, Description: "公开"},
			{Value: 1, Description: "粉丝"},
		}},
		sourceFieldFactory("内容声明", []OptionDef{
			{Value: 0, Description: "不声明"},
			{Value: 1, Description: "内容由AI生成"},
			{Value: 2, Description: "内容为虚构演绎"},
		}),
		{Name: "dynamic", Type: "string", Description: "粉丝动态"},
		timerPublishField,
	}
}

func baijiahaoBaseSetting() []FieldDef {
	return []FieldDef{
		{Name: "location", Type: "string", Description: "位置"},
		{Name: "classify", Type: "string", Description: "分类（格式：\"一级分类/二级分类\" 或 \"一级/二级/三级\"）"},
		{Name: "activity", Type: "string", Description: "活动"},
		{Name: "byAI", Type: "bool", Description: "AI创作声明"},
		timerPublishField,
	}
}
