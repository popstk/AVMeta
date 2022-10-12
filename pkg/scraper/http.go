package scraper

import (
	"crypto/tls"
	"github.com/ahuigo/requests"
	"net/http"
	"time"
)

const (
	DefaultUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"
)

func MapToCookie(m map[string]string) []*http.Cookie {
	ret := make([]*http.Cookie, 0, len(m))
	for k, v := range m {
		ret = append(ret, &http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	return ret
}

func RequestSession(cookies []*http.Cookie, ua string, retry int, timeout time.Duration, proxies string, verify bool) *requests.Session {
	s := requests.NewSession()
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			s.SetCookie(cookie)
		}
	}
	if len(ua) > 0 {
		s.SetGlobalHeader("User-Agent", ua)
	}

	if timeout > 0 {
		s.SetTimeout(timeout)
	}

	if len(proxies) > 0 {
		s.Proxy(proxies)
	}

	if !verify {
		tr := s.Client.Transport.(*http.Transport)
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: verify}
	}

	return s
}
