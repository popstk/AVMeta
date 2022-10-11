package scraper

import (
	"bytes"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"
)

// https://github.com/yoshiko2/Movie_Data_Capture/blob/2815762a8aecd7436da84d87a8546414a2dd8ef6/scrapinglib/javlibrary.py

// JavLibraryScraper jav library网站刮削器
type JavLibraryScraper struct {
	Site   string     // 免翻地址
	Proxy  string     // 代理配置
	uri    string     // 页面地址
	number string     // 最终番号
	root   *html.Node // 根节点
}

// NewJavLibraryScraper 返回一个被初始化的javbus刮削对象
//
// site 字符串参数，传入免翻地址，
// proxy 字符串参数，传入代理信息
func NewJavLibraryScraper(site, proxy string) *JavLibraryScraper {
	return &JavLibraryScraper{Site: site, Proxy: proxy}
}

// Fetch 刮削
func (s *JavLibraryScraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 获取信息
	return nil
}

// 搜索获取真正的url
func (s *JavLibraryScraper) queryNumberUrl() (string, error) {
	queryURL := "http://www.javlibrary.com/cn/vl_searchbyid.php?keyword=" + s.number
	client := util.NewProxyClient(s.Proxy)
	rsp, err := util.HttpGet(client, queryURL)
	if err != nil {
		return "", err
	}

	finalURL := rsp.Respond.Request.URL.String()

	node, err := htmlquery.Parse(bytes.NewReader(rsp.Body))
	if err != nil {
		return "", fmt.Errorf("parse body: %s, err: %w", string(rsp.Body), err)
	}

	if strings.Contains(finalURL, "/?v=jav") {
		s.root = node
		return finalURL, nil
	}

	numbers := FindListFromText(node, `//div[@class="id"]/text()`)
	found := -1
	for i, number := range numbers {
		if number == s.number {
			found = i
			break
		}
	}

	if found < 0 {
		return "", nil
	}

	urls := FindListFromText(node, `//div[@class="id"]/../@href`)
	if found >= len(urls) {
		return "", fmt.Errorf("can not find url with index %d, urls: %+v", found, urls)
	}

	return urls[found], nil
}

// GetTitle 获取名称
func (s *JavLibraryScraper) GetTitle() string {
	return ""
}

// GetIntro 获取简介
func (s *JavLibraryScraper) GetIntro() string {
	return GetDmmIntro(s.number, s.Proxy)
}

// GetDirector 获取导演
func (s *JavLibraryScraper) GetDirector() string {
	return ""
}

// GetRelease 发行时间
func (s *JavLibraryScraper) GetRelease() string {
	return ""
}

// GetRuntime 获取时长
func (s *JavLibraryScraper) GetRuntime() string {
	return ""
}

// GetStudio 获取厂商
func (s *JavLibraryScraper) GetStudio() string {
	return ""
}

// GetSeries 获取系列
func (s *JavLibraryScraper) GetSeries() string {
	return ""
}

// GetTags 获取标签
func (s *JavLibraryScraper) GetTags() []string {
	// 类别数组
	var tags []string
	return tags
}

// GetCover 获取图片
func (s *JavLibraryScraper) GetCover() string {
	return ""
}

// GetActors 获取演员
func (s *JavLibraryScraper) GetActors() map[string]string {
	// 演员数组
	actors := make(map[string]string)

	return actors
}

// GetURI 获取页面地址
func (s *JavLibraryScraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *JavLibraryScraper) GetNumber() string {
	return s.number
}
