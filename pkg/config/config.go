package config

import (
	"github.com/spf13/viper"
	"sort"
)

// BaseStruct 配置信息基础节点
type BaseStruct struct {
	Proxy string // 代理地址
}

// PathStruct 配置信息路径节点
type PathStruct struct {
	Success   string   // 成功存储目录
	Fail      string   // 失败存储目录
	Directory string   // 影片存储路径格式
	Filter    []string // 文件名过滤规则
}

// MediaStruct 配置信息媒体库节点
type MediaStruct struct {
	Library   string // 媒体库类型
	URL       string // Emby访问地址
	API       string // Emby API Key
	SecretID  string // 腾讯云 SecretId
	SecretKey string // 腾讯云 SecretKey
}

// SiteStruct 配置信息网站节点
type SiteStruct struct {
	JavBus string // javbus免翻地址
	JavDB  string // javdb免翻地址
}

// Conf 程序配置信息结构
type Conf struct {
	Base  BaseStruct  // 基础配置
	Path  PathStruct  // 路径配置
	Media MediaStruct // 媒体库配置
	Site  SiteStruct  // 免翻地址配置
	Code  []string    // 优先匹配番号

	IgnoreDir map[string]struct{} `yaml:"-"`
}

// GetConfig 读取配置信息，返回配置信息对象，
// 若没有配置文件，则创建一份默认配置文件并读取返回。
func GetConfig() (*Conf, error) {
	// 配置名称
	viper.SetConfigName("config")
	// 配置类型
	viper.SetConfigType("yaml")
	// 添加当前执行路径为配置路径
	viper.AddConfigPath(".")
	// 读取配置信息
	err := viper.ReadInConfig()
	// 读取配置
	if err != nil {
		// 如果文件不存在
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return WriteConfig()
		}

		// 直接返回错误信息
		return nil, err
	}

	// 定义配置变量
	var config Conf

	// 反序列
	err = viper.Unmarshal(&config)

	// 初始化配置
	sort.Slice(config.Path.Filter, func(i, j int) bool {
		if len(config.Path.Filter[i]) != len(config.Path.Filter[j]) {
			return len(config.Path.Filter[i]) > len(config.Path.Filter[j])
		}

		return config.Path.Filter[i] < config.Path.Filter[j]
	})

	return &config, err
}

// WriteConfig 在程序执行路径下写入一份默认配置文件，
// 若写入成功则返回配置信息，若写入失败，则返回错误信息。
func WriteConfig() (*Conf, error) {
	// 配置名称
	viper.SetConfigName("config")
	// 配置类型
	viper.SetConfigType("yaml")
	// 添加当前执行路径为配置路径
	viper.AddConfigPath(".")

	// 默认配置
	cfg := &Conf{
		Base: BaseStruct{
			Proxy: "",
		},
		Path: PathStruct{
			Success:   "success",
			Fail:      "fail",
			Directory: "{number}",
			Filter:    []string{"thz.la"},
		},
		Media: MediaStruct{
			Library:   "nfo",
			URL:       "",
			API:       "",
			SecretID:  "",
			SecretKey: "",
		},
		Site: SiteStruct{
			JavBus: "https://www.javbus.com/",
			JavDB:  "https://javdb4.com/",
		},
	}

	// 设置数据
	viper.Set("base", cfg.Base)
	viper.Set("path", cfg.Path)
	viper.Set("media", cfg.Media)
	viper.Set("site", cfg.Site)
	viper.Set("code", cfg.Code)

	return cfg, viper.SafeWriteConfig()
}

func (c *Conf) Init() {
	c.IgnoreDir = make(map[string]struct{})
	c.IgnoreDir[c.Path.Success] = struct{}{}
	c.IgnoreDir[c.Path.Fail] = struct{}{}
}

func (c *Conf) InIgnoreDir(dir string) bool {
	_, ok := c.IgnoreDir[dir]
	return ok
}
