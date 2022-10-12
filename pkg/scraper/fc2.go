package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/ylqjgm/AVMeta/pkg/util"
)

const (
	exprTitle       = `/html/head/title/text()`
	exprStudio      = `//*[@id="top"]/div[1]/section[1]/div/section/div[2]/ul/li[3]/a/text()`
	exprRelease     = `//*[@id="top"]/div[1]/section[1]/div/section/div[2]/div[2]/p/text()`
	exprRuntime     = `//p[@class='items_article_info']/text()`
	exprDirector    = `//*[@id="top"]/div[1]/section[1]/div/section/div[2]/ul/li[3]/a/text()`
	exprActor       = `//*[@id="top"]/div[1]/section[1]/div/section/div[2]/ul/li[3]/a/text()`
	exprCover       = `//div[@class='items_article_MainitemThumb']/span/img/@src`
	exprExtraFanArt = `//ul[@class="items_article_SampleImagesArea"]/li/a/@href`
	exprTags        = `//a[@class='tag tagTag']/text()`
)

// FC2Scraper fc2网站刮削器
type FC2Scraper struct {
	Proxy   string     // 代理设置
	uri     string     // 页面地址
	code    string     // 临时番号
	number  string     // 最终番号
	fc2Root *html.Node // fc2根节点
}

// fc2标签json结构
type fc2tags struct {
	Tags []fc2tag `json:"tags"`
}

// fc2标签内容结构
type fc2tag struct {
	Tag string `json:"tag"`
}

// NewFC2Scraper 返回一个被初始化的fc2刮削对象
//
// proxy 字符串参数，传入代理信息
func NewFC2Scraper(proxy string) *FC2Scraper {
	return &FC2Scraper{Proxy: proxy}
}

// Fetch 刮削
func (s *FC2Scraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 过滤番号
	r := regexp.MustCompile(`\d{6,7}`)
	// 获取临时番号
	s.code = r.FindString(code)
	// 组合fc2地址
	fc2uri := fmt.Sprintf("https://adult.contents.fc2.com/article/%s/", s.code)

	data, status, err := util.MakeRequest("GET", fc2uri, s.Proxy, nil, nil, nil)
	if err != nil || status >= http.StatusBadRequest {
		return fmt.Errorf("%s [fetch]: status: %d", fc2uri, status)
	}

	fc2Root, err := htmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		return err
	}

	s.uri = fc2uri
	s.fc2Root = fc2Root
	return nil
}

func FindFromText(r *html.Node, expr string) string {
	node := htmlquery.FindOne(r, expr)
	if node == nil {
		return ""
	}

	return htmlquery.InnerText(node)
}

func FindListFromText(r *html.Node, expr string) []string {
	nodes := htmlquery.Find(r, expr)
	if len(nodes) == 0 {
		return nil
	}
	result := make([]string, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, htmlquery.InnerText(node))
	}

	return result
}

// GetTitle 获取名称
func (s *FC2Scraper) GetTitle() string {
	return FindFromText(s.fc2Root, exprTitle)
}

// GetOutline 获取简介
func (s *FC2Scraper) GetOutline() string {
	return ""
}

// GetDirector 获取导演
func (s *FC2Scraper) GetDirector() string {
	return FindFromText(s.fc2Root, exprDirector)
}

// GetRelease 发行时间
func (s *FC2Scraper) GetRelease() string {
	val := FindFromText(s.fc2Root, exprRelease)
	return strings.ReplaceAll(strings.Trim(val, " ['販売日 : ']"), "/", "-")
}

// GetRuntime 获取时长
func (s *FC2Scraper) GetRuntime() string {
	return "0"
}

// GetStudio 获取厂商
func (s *FC2Scraper) GetStudio() string {
	return util.FC2
}

// GetSeries 获取系列
func (s *FC2Scraper) GetSeries() string {
	return util.FC2
}

// GetTags 获取标签
func (s *FC2Scraper) GetTags() []string {
	// 组合地址
	uri := fmt.Sprintf("https://adult.contents.fc2.com/api/v4/article/%s/tag?", s.code)

	// 读取远程数据
	data, err := util.GetResult(uri, s.Proxy, nil)
	// 检查
	if err != nil {
		return nil
	}

	// 读取内容
	body, err := ioutil.ReadAll(bytes.NewReader(data))
	// 检查错误
	if err != nil {
		return nil
	}

	// json
	var tagsJSON fc2tags

	// 解析json
	err = json.Unmarshal(body, &tagsJSON)
	// 检查
	if err != nil {
		return nil
	}

	// 定义数组
	var tags []string

	// 循环标签
	for _, tag := range tagsJSON.Tags {
		tags = append(tags, strings.TrimSpace(tag.Tag))
	}

	return tags
}

// GetCover 获取图片
func (s *FC2Scraper) GetCover() string {
	val := FindFromText(s.fc2Root, exprCover)
	if len(val) == 0 {
		return ""
	}

	p, _ := util.JoinPath("https://adult.contents.fc2.com", val)
	return p
}

// GetActors 获取演员
func (s *FC2Scraper) GetActors() map[string]string {
	node := htmlquery.FindOne(s.fc2Root, exprActor)
	if node == nil {
		return map[string]string{
			"素人": "",
		}
	}

	val := htmlquery.InnerText(node)
	return map[string]string{
		val: "",
	}
}

// GetURI 获取页面地址
func (s *FC2Scraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *FC2Scraper) GetNumber() string {
	return s.number
}
