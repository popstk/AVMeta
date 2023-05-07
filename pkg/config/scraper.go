package config

type Scraper struct {
	Disable    bool   `yaml:"disable"`
	Site       string `yaml:"site"`
	UseProxy   bool   `yaml:"use_proxy"`
	Proxy      string `yaml:"proxy"`
	CookieFile string `yaml:"cookie_file"`
}

func (s *Scraper) GetProxy() string {
	if !s.UseProxy {
		return ""
	}

	return s.Proxy
}

func (s *Scraper) SetDefaultProxy(p string) {
	if len(s.Proxy) == 0 {
		s.Proxy = p
	}
}
