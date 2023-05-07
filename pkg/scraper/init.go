package scraper

import "github.com/ylqjgm/AVMeta/pkg/config"

var RegisterScraper = make(map[string]func(config.Scraper) IScraper)
