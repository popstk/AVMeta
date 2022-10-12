package scraper

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/ylqjgm/AVMeta/pkg/util"
	"golang.org/x/net/html"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SiroScraper siro网站刮削器
type SiroScraper struct {
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
	node   *html.Node
}

// NewSiroScraper 返回一个被初始化的siro刮削对象
//
// proxy 字符串参数，传入代理信息
func NewSiroScraper(proxy string) *SiroScraper {
	return &SiroScraper{Proxy: proxy}
}

// Fetch 刮削
func (s *SiroScraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)

	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:   "adc",
		Value:  "1",
		Path:   "/",
		Domain: ".mgstage.com",
	})

	uri := fmt.Sprintf("https://www.mgstage.com/product/product_detail/%s/", s.number)

	// 执行请求
	data, status, err := util.MakeRequest("GET", uri, s.Proxy, nil, nil, cookies)
	if err != nil {
		return err
	}

	if status >= http.StatusBadRequest {
		return fmt.Errorf("%s [Http Status]: %d", uri, status)
	}

	root, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return err
	}

	// 设置页面地址
	s.uri = uri
	// 设置根节点
	s.root = root

	return nil
}

// GetTitle 获取名称
func (s *SiroScraper) GetTitle() string {
	return s.root.Find(`h1.tag`).Text()
}

// GetOutline 获取简介
func (s *SiroScraper) GetOutline() string {
	return util.IntroFilter(s.root.Find(`#introduction p.introduction`).Text())
}

// GetDirector 获取导演
func (s *SiroScraper) GetDirector() string {
	return ""
}

// GetRelease 发行时间
func (s *SiroScraper) GetRelease() string {
	return s.root.Find(`th:contains("配信開始日")`).NextFiltered("td").Text()
}

// GetRuntime 获取时长
func (s *SiroScraper) GetRuntime() string {
	return strings.TrimRight(s.root.Find(`th:contains("収録時間")`).NextFiltered("td").Text(), "min")
}

// GetStudio 获取厂商
func (s *SiroScraper) GetStudio() string {
	val := s.root.Find(`th:contains("メーカー")`).NextFiltered("td").Text()
	if val == "" {
		val = s.root.Find(`th:contains("メーカー")`).NextFiltered("td").Find("a").Text()
	}

	return val
}

// GetSeries 获取系列
func (s *SiroScraper) GetSeries() string {
	val := s.root.Find(`th:contains("シリーズ")`).NextFiltered("td").Text()
	if val == "" {
		val = s.root.Find(`th:contains("シリーズ")`).NextFiltered("td").Find("a").Text()
	}

	return val
}

// GetTags 获取标签
func (s *SiroScraper) GetTags() []string {
	// 标签数组
	var tags []string
	// 循环获取
	s.root.Find(`th:contains("ジャンル")`).NextFiltered("td").Find("a").Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetCover 获取图片
func (s *SiroScraper) GetCover() string {
	// 获取图片
	cover, exist := s.root.Find(`#EnlargeImage`).Attr("href")
	if !exist {
		ret, err := s.root.Html()
		log.Errorf("root err: %v text: %s", err, ret)
	}

	return cover
}

// GetActors 获取演员
func (s *SiroScraper) GetActors() map[string]string {
	// 演员数组
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`th:contains("出演")`).NextFiltered("td").Find("a").Each(func(i int, item *goquery.Selection) {
		// 演员名字
		actors[strings.TrimSpace(item.Text())] = ""
	})

	// 是否获取到
	if len(actors) == 0 {
		// 重新获取
		name := s.root.Find(`th:contains("出演")`).NextFiltered("td").Text()
		// 获取演员名字
		actors[strings.TrimSpace(name)] = ""
	}

	return actors
}

// GetURI 获取页面地址
func (s *SiroScraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *SiroScraper) GetNumber() string {
	return s.number
}
