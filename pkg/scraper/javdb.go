package scraper

import (
	"bytes"
	"fmt"
	"github.com/ahuigo/requests"
	log "github.com/sirupsen/logrus"
	"github.com/ylqjgm/AVMeta/pkg/config"
	"net/http"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

const JavDB = "javdb"

func init() {
	RegisterScraper[JavDB] = func(c config.Scraper) IScraper {
		return NewJavDBScraper(c)
	}
}

// JavDBScraper javdb网站刮削器
type JavDBScraper struct {
	conf *config.Scraper

	session *requests.Session
	uri     string            // 页面地址
	number  string            // 最终番号
	root    *goquery.Document // 根节点
}

// NewJavDBScraper 返回一个被初始化的javdb刮削对象
//
// site 字符串参数，传入免翻地址，
// proxy 字符串参数，传入代理信息
func NewJavDBScraper(conf config.Scraper) *JavDBScraper {
	return &JavDBScraper{conf: &conf}
}

// Fetch 刮削
func (s *JavDBScraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)

	var cookies []*http.Cookie
	if len(s.conf.CookieFile) > 0 {
		val, err := ReadCookieFromFile(s.conf.CookieFile)
		if err != nil {
			return err
		}
		cookies = val
	}

	s.session = RequestSession(cookies, DefaultUA, 3, 0, s.conf.GetProxy(), false)

	// 搜索
	id, err := s.search()
	// 检查错误
	if err != nil {
		return fmt.Errorf("%s [Search]: %s", code, err)
	}

	// 组合地址
	uri := fmt.Sprintf("%s%s", util.CheckDomainPrefix(s.conf.Site), id)

	log.Infof("javdb found real uri: %s", uri)
	rsp, err := s.session.Get(uri)
	if err != nil {
		return fmt.Errorf("%s [fetch]: %s", uri, err)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewReader(rsp.Body()))
	if err != nil {
		return fmt.Errorf("%s [goquery]: %s", uri, err)
	}

	// 设置页面地址
	s.uri = uri
	// 设置根节点
	s.root = root

	return nil
}

// 搜索影片
func (s *JavDBScraper) search() (string, error) {
	// 组合地址
	uri := fmt.Sprintf("%s/search?q=%s&f=all", util.CheckDomainPrefix(s.conf.Site), strings.ToUpper(s.number))

	rsp, err := s.session.Get(uri)
	if err != nil {
		return "", err
	}

	content := rsp.Body()
	log.Debugf("content: %s", string(content))

	root, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return "", err
	}

	// 查找是否获取到
	if -1 < root.Find(`.empty-message:contains("暫無內容")`).Index() {
		return "", fmt.Errorf("404 Not Found")
	}

	// 定义ID
	var id string

	// 循环检查番号
	items := root.Find(`.movie-list>.item>a`)
	log.Debugf("items: %s", items.Text())

	items.Each(func(i int, item *goquery.Selection) {
		// 获取番号
		date := item.Find(".video-title>strong").Text()
		// 大写并去除空白
		date = strings.ToUpper(strings.TrimSpace(date))
		// 检查番号是否完全正确
		if strings.EqualFold(strings.ToUpper(s.number), date) {
			// 获取href元素
			id, _ = item.Attr("href")
		}
	})

	// 清除空白
	id = strings.TrimSpace(id)

	// 是否获取到
	if id == "" {
		return "", fmt.Errorf("%s [fetch]: 404 Not Found ID", uri)
	}

	return id, nil
}

// GetTitle 获取名称
func (s *JavDBScraper) GetTitle() string {
	return s.root.Find(`div.video-detail .current-title`).Text()
}

// GetOutline 获取简介
func (s *JavDBScraper) GetOutline() string {
	return GetDmmIntro(s.number, s.conf.GetProxy())
}

// GetDirector 获取导演
func (s *JavDBScraper) GetDirector() string {
	// 获取数据
	val := s.root.Find(`strong:contains("導演")`).NextFiltered(`span.value`).Text()
	// 检查
	if val == "" {
		val = s.root.Find(`strong:contains("導演")`).NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetRelease 发行时间
func (s *JavDBScraper) GetRelease() string {
	// 获取数据
	val := s.root.Find(`strong:contains("日期")`).NextFiltered(`span.value`).Text()
	// 检查
	if val == "" {
		val = s.root.Find(`strong:contains("日期")`).NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetRuntime 获取时长
func (s *JavDBScraper) GetRuntime() string {
	// 获取数据
	val := s.root.Find(`strong:contains("時長")`).NextFiltered(`span.value`).Text()
	// 检查
	if val == "" {
		val = s.root.Find(`strong:contains("時長")`).NextFiltered(`span.value`).Find("a").Text()
	}

	// 去除多余
	val = strings.TrimRight(val, "分鍾")

	return val
}

// GetStudio 获取厂商
func (s *JavDBScraper) GetStudio() string {
	// 获取数据
	val := s.root.Find(`strong:contains("片商")`).NextFiltered(`span.value`).Text()
	// 检查
	if val == "" {
		val = s.root.Find(`strong:contains("片商")`).NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetSeries 获取系列
func (s *JavDBScraper) GetSeries() string {
	// 获取数据
	val := s.root.Find(`strong:contains("系列")`).NextFiltered(`span.value`).Text()
	// 检查
	if val == "" {
		val = s.root.Find(`strong:contains("系列")`).NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetTags 获取标签
func (s *JavDBScraper) GetTags() []string {
	// 类别数组
	var tags []string
	// 循环获取
	s.root.Find(`strong:contains("類別")`).NextFiltered(`span.value`).Find("a").Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetCover 获取图片
func (s *JavDBScraper) GetCover() string {
	// 获取图片
	fanart, _ := s.root.Find(`div.column-video-cover a img`).Attr("src")

	return fanart
}

// GetActors 获取演员
func (s *JavDBScraper) GetActors() map[string]string {
	// 演员列表
	actors := make(map[string]string)

	node := s.root.Find(`strong:contains("演員")`).NextFiltered(`span.value`)
	var lastActor string
	node.Children().Each(func(i int, selection *goquery.Selection) {
		if i%2 == 0 {
			lastActor = selection.Text()
			return
		}
		if selection.HasClass("female") {
			actors[lastActor] = ""
		}
	})

	return actors
}

// GetURI 获取页面地址
func (s *JavDBScraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *JavDBScraper) GetNumber() string {
	return s.number
}
