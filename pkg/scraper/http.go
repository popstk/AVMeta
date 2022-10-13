package scraper

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ahuigo/requests"
	"math"
	"net/http"
	"os"
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

type CookieJson struct {
	Domain         string      `json:"domain"`
	ExpirationDate float64     `json:"expirationDate"`
	HostOnly       bool        `json:"hostOnly"`
	HttpOnly       bool        `json:"httpOnly"`
	Name           string      `json:"name"`
	Path           string      `json:"path"`
	SameSite       string      `json:"sameSite"`
	Secure         bool        `json:"secure"`
	Session        bool        `json:"session"`
	StoreId        interface{} `json:"storeId"`
	Value          string      `json:"value"`
}

func ReadCookieFromFile(file string) ([]*http.Cookie, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file %s err: %w", file, err)
	}
	var cookies []CookieJson
	if err := json.Unmarshal(data, &cookies); err != nil {
		return nil, fmt.Errorf("json.Unmarshal err: %w", err)
	}

	ret := make([]*http.Cookie, 0, len(cookies))
	for _, c := range cookies {
		sec, dec := math.Modf(c.ExpirationDate)

		ret = append(ret, &http.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Path:     c.Path,
			Domain:   c.Domain,
			Expires:  time.Unix(int64(sec), int64(dec*(1e9))),
			Secure:   c.Secure,
			HttpOnly: c.HttpOnly,
		})
	}

	return ret, nil
}
