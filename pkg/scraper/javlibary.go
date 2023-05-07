package scraper

import (
	"bytes"
	"fmt"
	"github.com/ahuigo/requests"
	"github.com/antchfx/htmlquery"
	log "github.com/sirupsen/logrus"
	"strings"
)

// https://github.com/yoshiko2/Movie_Data_Capture/blob/2815762a8aecd7436da84d87a8546414a2dd8ef6/scrapinglib/javlibrary.py

var JavLibraryExpr = ParserExpr{
	Number:      `//div[@id="video_id"]/table/tr/td[@class="text"]/text()`,
	Title:       `//div[@id="video_title"]/h3/a/text()`,
	Studio:      `//div[@id="video_maker"]/table/tr/td[@class="text"]/span/a/text()`,
	Release:     `//div[@id="video_date"]/table/tr/td[@class="text"]/text()`,
	Runtime:     `//div[@id="video_length"]/table/tr/td/span[@class="text"]/text()`,
	Director:    `//div[@id="video_director"]/table/tr/td[@class="text"]/span/a/text()`,
	Actor:       `//div[@id="video_cast"]/table/tr/td[@class="text"]/span/span[@class="star"]/a/text()`,
	Cover:       `//img[@id="video_jacket_img"]/@src`,
	ExtraFanArt: `//div[@class="previewthumbs"]/img/@src`,
	Tags:        `//div[@id="video_genres"]/table/tr/td[@class="text"]/span/a/text()`,
	UserRating:  `//div[@id="video_review"]/table/tr/td/span[@class="score"]/text()`,
}

// JavLibraryScraper jav library网站刮削器
type JavLibraryScraper struct {
	Site  string // 免翻地址
	Proxy string // 代理配置

	session *requests.Session
	Scraper
}

func NewJavLibraryScraper(proxy string) *JavLibraryScraper {
	return &JavLibraryScraper{
		Proxy:   proxy,
		Scraper: Scraper{Expr: JavLibraryExpr},
	}
}

// Fetch 刮削
func (s *JavLibraryScraper) Fetch(code string) error {
	keyValue := make(map[string]string)
	keyValue["over18"] = "1"

	var cookies []string
	for _, value := range cookies {
		p := strings.SplitN(value, "=", 2)
		if len(p) != 2 {
			return fmt.Errorf("invalid cookie %s", value)
		}
		keyValue[p[0]] = p[1]
	}

	// 设置番号
	s.number = strings.ToUpper(code)
	s.session = RequestSession(MapToCookie(keyValue), DefaultUA, 3, 0, s.Proxy, false)

	var err error
	s.uri, err = s.queryNumberUrl()
	if err != nil {
		return fmt.Errorf("queryNumberUrl: %w", err)
	}

	log.Infof("queryNumberUrl uri %s", s.uri)

	if s.root == nil {
		rsp, err := s.session.Get(s.uri)
		if err != nil {
			return fmt.Errorf("session get: %w", err)
		}
		s.root, err = htmlquery.Parse(bytes.NewReader(rsp.Body()))
		if err != nil {
			return fmt.Errorf("parse body: %s, err: %w", rsp.Body(), err)
		}
	}

	return nil
}

// 搜索获取真正的url
func (s *JavLibraryScraper) queryNumberUrl() (string, error) {
	queryURL := "http://www.javlibrary.com/cn/vl_searchbyid.php?keyword=" + s.number

	rsp, err := s.session.Get(queryURL)
	if err != nil {
		return "", err
	}

	finalURL := rsp.R.Request.URL.String()
	log.Infof("finalURL %s", finalURL)

	node, err := htmlquery.Parse(bytes.NewReader(rsp.Body()))
	if err != nil {
		return "", fmt.Errorf("parse body: %s, err: %w", rsp.Body(), err)
	}

	log.Infof("html: %s", htmlquery.OutputHTML(node, false))

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

	log.Infof("found index %d", found)

	urls := FindListFromText(node, `//div[@class="id"]/../@href`)
	if found >= len(urls) {
		return "", fmt.Errorf("can not find url with index %d, urls: %+v", found, urls)
	}

	return urls[found], nil
}

func (s *JavLibraryScraper) GetTitle() string {
	title := s.Scraper.GetTitle()
	return strings.TrimSpace(strings.ReplaceAll(title, s.GetNumber(), ""))
}

func (s *JavLibraryScraper) GetCover() string {
	uri := s.Scraper.GetCover()
	if !strings.HasPrefix(uri, "http") {
		uri = "https" + uri
	}

	return uri
}

func (s *JavLibraryScraper) GetOutline() string {
	return ""
}
