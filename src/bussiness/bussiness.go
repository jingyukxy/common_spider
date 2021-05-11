package bussiness

import (
	"time"
)

// 应用类型
type AppType struct {
	Id          uint64 // id
	Name        string // 名称
	Pid         uint64 // 父Id
	Status      bool   // 状态
	Description string // 描述
	Keywords    string // 关键字
}

// 应用
type App struct {
	Id          uint64    // ID
	Name        string    // 名称
	Title       string    // 标题
	AppType     AppType   // 应用类型
	Keywords    string    // 关键字
	HasMobile   bool      // 最否有移动端
	HasSearch   bool      // 是否有搜索
	Status      bool      // 状态
	AppUrls     []AppUrl  // 关联链接
	UpdateTime  time.Time // 更新时间
	CreateTime  time.Time // 创建时间
	Description string    // 描述
}

// api类型
type ApiType struct {
	Id          uint64 // ID
	Pid         uint64 // 父ID
	TypeName    string // 类型名称
	Description string // 描述
}

// 应用URL
type AppUrl struct {
	Id         uint64    // id
	Name       string    // 名称
	Url        string    // 链接
	ApiType    ApiType   // 对应api类型
	DeviceType uint8     // 设备类型
	Status     bool      // 状态
	Theme      string    // 对应主题
	CreateTime time.Time // 创建时间
}

// 广告
type Ads struct {
	Id             uint64    // id
	AdsMerchant    uint64    // 厂商id
	AdsPosition    uint64    // 位置id
	AdsTemplate    uint64    // 模板id
	AdsCostType    uint64    //计费方式
	AdsType        uint8     // 类型
	Image          string    // 图片
	ImgHeight      int       // 图片高度
	ImgWidth       int       // 图片宽度
	IsFlash        bool      // 是否为flash
	MerchantRemark string    // 厂商描述
	OpenMode       uint8     // 打开方式
	Weight         int       // 权重
	OtherSetting   string    // 其它设置
	Status         bool      // 状态
	Expiration     time.Time // 过期时间
	StartTime      time.Time // 开始时间
	CreateTime     time.Time // 创建时间
}

// 广告主
type AdsMerchant struct {
	Id          uint64    // id
	Name        string    // 名称
	Company     string    // 公司
	Phone       string    // 电话
	QQ          string    // qq
	Description string    // 描述
	CreateTime  time.Time // 创建时间
}

// 广告版位
type AdsPosition struct {
	Id           uint64    // id
	Name         string    // 名称
	Description  string    // 描述
	AppId        uint64    // APPId
	AppUrl       string    // APP对应的url
	Height       int       // 高度
	Width        int       // 宽度
	Status       bool      // 状态
	DisplayMode  uint8     // 展示方式
	CarouselTime int       // 轮播时间
	CreateTime   time.Time // 创建时间
	UpdateTime   time.Time // 更新时间
}

// 广告模板
type AdsTemplate struct {
	Id         uint64    // id
	Name       string    // 名称
	JSCode     string    // js代码
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 更新时间
	AdsMode    uint8     // 广告模式
}

// 计费模式
type AdsCost struct {
	Id          uint64  // id
	Name        string  // 名称
	CostType    uint8   // 计费类型 cmp cpc cpt cpa ...
	Price       float64 // 单价
	Description string  // 描述
}

// 外链
type ExternalLink struct {
	Id         uint64    // id
	Name       string    // 链接名称
	AppId      uint64    // APPid
	AppUrl     string    // APPurl
	LinkUrl    string    // 外链地址
	LinkType   int8      // 链接类型
	Logo       string    // logo地址
	StartTime  time.Time // 开始时间
	Expiration time.Time // 过期时间
}

// 用户表
type User struct {
	Id           uint64    // id
	UserName     string    // 用户名
	Phone        string    // 电话
	Email        string    // 邮箱
	Password     string    // 密码
	Avatar       string    // 头像
	AppId        uint64    // Appid
	Salt         string    // 密钥
	LoginTimes   uint64    // 登录次数
	IsThirdParty bool      // 是否为第三方用户
	OpenType     int8      // 第三方类型
	OpenId       string    // 第三方openid
	LastLogin    time.Time // 最后登录时间
	Status       bool      //状态
	CreateTime   time.Time //创建时间
	UpdateTime   time.Time // 更新时间
}

// 用户信息表
type UserProfile struct {
	Id     uint64    // id
	Uid    uint64    // uid
	Birth  time.Time // 生日
	Gender int8      // 性别
	Level  int       // 级别
	Score  int       // 积分
	IsVip  bool      // 是否为VIP
}

// DNS厂商
type DNSMerchant struct {
	Id           uint64 // id
	Name         string // 名称
	MerchantType int8   // 厂商类型
	Description  string // 描述
	HasApi       bool   // 是否有API
}

// 域名
type Domain struct {
	Id         uint64    // domain id
	Dns        uint64    // dns厂商id
	Main       string    // 主域名
	Expiration time.Time // 过期时间
	IsUsed     bool      // 是否在用
}

// 子域名
type SubDomain struct {
	Id         uint64 // id
	DomainId   uint64 // 域名id
	SubName    string // 子域名 a b c
	DomainType int8   //  域名类型 txt a cname mail...
	SubDest    string // 指向目标
	TTL        int    // ttl
}

// 采集器
type Spider struct {
	Id            uint64    // id
	Name          string    // 名称
	Description   string    // 描述
	AppType       uint64    // app类型
	SourceIp      string    // 源ip
	QueueName     string    // 队列名称
	Rule          string    // rule规则 ajax/html
	Status        string    // 状态
	NeedDownLoad  bool      // 是否需要下载
	NeedTranslate bool      // 是否需要翻译
	NeedTag       bool      // 是否需要打标签
	NeedProxy     bool      // 是否需要代理
	CreateTime    time.Time // 创建时间
}

// 采集任务
type SpiderTask struct {
	Id           uint64    // id
	TaskName     string    // 任务名称
	Description  string    // 描述
	CreateTime   time.Time // 创建时间
	Status       bool      // 状态
	TaskUrl      string    // 任务url
	TaskSchedule string    // 任务调度设定
}

// 任务日志
type SpiderTaskLog struct {
	Id         uint64    // id
	TaskId     uint64    // 任务id
	ExecTime   time.Time // 执行时间
	ExecStatus int8      // 执行状态，执行中 执行成功 执行失败
	ExecLog    string    // 执行日志
}

// 代理配置
type ProxyConfig struct {
	Id         uint64
	Addr       string
	Protocol   string
	CreateTime time.Time
	Status     bool
}

// 第三方Api配置
type ThirdApiConfig struct {
	Id            uint64    // id
	Name          string    // 名称
	ApiAddr       string    //  api地址
	ApiType       string    // api类型
	AuthToken     string    // 认证token
	CreateTime    time.Time // 创建时间
	Status        bool      // 状态
	OtherSettings string    // 其它设定
}

//标签
type Tag struct {
	Id      uint64 // id
	AppType uint64 // apptype id
	Weight  int    // 权重
	Name    string // 名称
	Part    string // 词性
}

// 同义词
type Synonym struct {
	Id      uint64 // id
	Word    string // 名称
	SynWord string // 同义词
}

// 菜单
type Menu struct {
	Id          uint64    // id
	AppId       uint64    // app id
	AppType     uint64    // app类型
	AppUrl      string    // app链接
	Name        string    // 名称
	Pid         uint64    // pid
	Level       int       // 层级
	Status      bool      // 状态
	Component   string    // 组件
	ExternalUrl string    // 外部连接
	Icon        string    // 图标
	Title       string    // 标题
	Keywords    string    // 关键字
	Description string    //描述
	CreateTime  time.Time // 创建时间
}

// 评论
type Comment struct {
	Id         uint64    // id
	AppId      uint64    // app id
	AppType    uint64    // app type id
	CategoryId uint64    // 分类Id
	ContentId  uint64    // 内容Id
	Title      string    // 分类标题
	UserName   string    // 评论用户
	Ip         string    // ip
	CreateTime time.Time // 评论时间
	AuditUser  string    // 审核用户
	Status     int8      // 状态
	Content    string    // 评论内容
}

// ---------------------------------------------------------
// 栏目
type Category struct {
	Id          uint64 // id
	Name        string // 名称
	Pid         uint64 // pid
	Level       int    // 层级
	Status      bool   // 状态
	Keywords    string // 关键字
	Title       string // 标题
	Description string // 描述
	CatType     int8   // 类型 1.普通 2.专题...
	AppId       uint64 // App id
	AppName     string // app name
	AppType     uint64 // app type
}

// 文章
type Article struct {
	Id             uint64    // id
	CategoryId     uint64    // 栏目id
	Title          string    // 标题
	SubTitle       string    // 副标题
	ThumbImg       string    // 缩略图
	SourceUrl      string    // 来源url
	Author         string    // 作者
	Tags           string    // 标签
	Description    string    // 描述
	CreateTime     time.Time // 创建时间
	PublishTime    time.Time // 发布时间
	UpdateTime     time.Time // 更新时间
	Status         int8      // 状态 1.保存 2.发布 3.删除
	IsTop          bool      // 是否置顶
	ReadDayCount   uint64    // 日阅读量
	ReadWeekCount  uint64    // 周阅读量
	ReadMonthCount uint64    // 月阅读量
}

// 文章内容
type ArticleContent struct {
	Id           uint64 // id
	ArticleId    uint64 // 文章id
	Content      string // 文章内容
	ContentType  int8   // 内容类型 文字/图片/文字图片
	AllowComment bool   //是否允许评论
	UpCount      uint64 // 点赞量
	DownCount    uint64 // 反对量
}

//----------------------------视频-----------------
// 演职人员
type VodStaff struct {
	Id          uint64    // Id
	Name        string    // 名称
	StaffType   int8      // 类型，1.演员 2.导演 3.编剧...
	Gender      int8      // 性别
	Image       string    // 图片
	Height      int8      // 身高
	Country     string    // 国家
	Description string    // 描述
	UpdateTime  time.Time // 更新时间
}

// 分集剧情
type VodEpisodes struct {
	Id          uint64    // id
	VodId       uint64    // 视频id
	Title       string    // 标题
	Description string    // 描述
	Keywords    string    // 关键字
	CreateTime  time.Time // 创建时间
}

// 视频
type Vod struct {
	Id            uint64    // id
	Name          string    // 名称
	CategoryId    uint64    // 类型
	Title         string    // 副标题
	Keywords      string    // 关键字
	Image         string    // 图片
	Area          string    // 地域
	Language      string    // 语言
	ShowTime      time.Time // 上映时间
	TotalEpisode  int       // 总集数
	NewEpisode    int       // 更新到
	IsEnd         bool      // 是否已完结
	ShowYear      int       // 上映年
	AddTime       time.Time // 添加时间
	HitDayCount   uint64    // 日点击
	HitWeekCount  uint64    // 周点击
	HitMonthCount uint64    // 月点击
	HitLastTime   time.Time // 最后点击时间
	Stars         int       // 星数
	Status        bool      // 状态
	Editor        string    // 编辑
	Score         float64   // 评分
	DoubanScore   float64   // 豆瓣评分
	IsFilm        bool      // 是否是电影
	VodLength     uint64    // 视频长度
	Description   string    // 描述
	UpdateTime    time.Time // 更新时间
}

// 视频标签
type VodTag struct {
	Id         uint64 // id
	Name       string // 名称
	CategoryId uint64 // 分类
}

// 视频内容
type VodContent struct {
	Id       uint64 // id
	Title    string // 标题
	PlayType string // 播放类型
	PlayUrl  string // 播放url
}

// 轮播
type VodCarousel struct {
	Id          uint64 // id
	Name        string // 名称
	Description string // 描述
	VodId       uint64 // 视频id
	CategoryId  uint64 // 分类id
	Status      bool   // 状态
	Image       string // 图片
	ExternalUrl string // 外部链接
}
