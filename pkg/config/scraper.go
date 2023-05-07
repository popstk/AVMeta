package config

type Scraper struct {
	Disable    bool   `mapstructure:"disable"`
	Site       string `mapstructure:"site"`
	UseProxy   bool   `mapstructure:"use_proxy"`
	Proxy      string `mapstructure:"proxy"`
	CookieFile string `mapstructure:"cookie_file"`
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
