package scraper

import "golang.org/x/net/html"

// IScraper 刮削器接口
type IScraper interface {
	// Fetch 执行刮削，并返回刮削结果
	//
	// code 字符串参数，传入番号信息
	Fetch(code string) error

	// GetURI 获取刮削的页面地址
	GetURI() string

	// GetNumber 获取最终的正确番号信息
	GetNumber() string

	// GetTitle 从刮削结果中获取影片标题
	GetTitle() string
	// GetOutline 从刮削结果中获取影片简介
	GetOutline() string
	// GetDirector 从刮削结果中获取影片导演
	GetDirector() string
	// GetRelease 从刮削结果中获取发行时间
	GetRelease() string
	// GetRuntime 从刮削结果中获取影片时长
	GetRuntime() string
	// GetStudio 从刮削结果中获取影片厂商
	GetStudio() string
	// GetSeries 从刮削结果中获取影片系列
	GetSeries() string
	// GetTags 从刮削结果中获取影片标签
	GetTags() []string
	// GetCover 从刮削结果中获取背景图片
	GetCover() string
	// GetActors 从刮削结果中获取影片演员
	GetActors() map[string]string
}

type ParserExpr struct {
	Number      string
	Title       string
	Studio      string
	Release     string
	Runtime     string
	Director    string
	Actor       string
	Cover       string
	ExtraFanArt string
	Tags        string
	UserRating  string
	Outline     string
	Series      string
}

type Scraper struct {
	Expr          ParserExpr
	Uncensored    bool
	MoreStoryLine bool
	SpecifiedUrl  string

	uri    string // 页面地址
	number string // 最终番号
	root   *html.Node
}

func (s *Scraper) Fetch(code string) error {
	return nil
}

func (s *Scraper) GetURI() string {
	return s.uri
}

func (s *Scraper) GetNumber() string {
	return s.number
}

func (s *Scraper) GetTitle() string {
	return FindFromText(s.root, s.Expr.Title)
}

func (s *Scraper) GetOutline() string {
	return FindFromText(s.root, s.Expr.Outline)
}

func (s *Scraper) GetDirector() string {
	return FindFromText(s.root, s.Expr.Director)
}

func (s *Scraper) GetRelease() string {
	return FindFromText(s.root, s.Expr.Release)
}

func (s *Scraper) GetRuntime() string {
	return FindFromText(s.root, s.Expr.Runtime)
}

func (s *Scraper) GetStudio() string {
	return FindFromText(s.root, s.Expr.Studio)
}

func (s *Scraper) GetSeries() string {
	return FindFromText(s.root, s.Expr.Series)
}

func (s *Scraper) GetTags() string {
	return FindFromText(s.root, s.Expr.Tags)
}

func (s *Scraper) GetCover() string {
	return FindFromText(s.root, s.Expr.Cover)
}

func (s *Scraper) GetActors() map[string]string {
	list := FindListFromText(s.root, s.Expr.Actor)
	m := make(map[string]string, len(list))
	for _, v := range list {
		m[v] = ""
	}

	return m
}
