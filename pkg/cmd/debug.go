package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ylqjgm/AVMeta/pkg/scraper"
)

func (e *Executor) initDebug() {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "调试",
		Long:  `调试模式`,
		Run:   e.debug,
	}
	cmd.Flags().StringP("proxy", "x", "", "代理")
	cmd.Flags().StringP("scraper", "s", "", "爬虫")
	e.rootCmd.AddCommand(cmd)
}

func (e *Executor) debug(c *cobra.Command, args []string) {
	proxy, err := c.Flags().GetString("proxy")
	if err != nil {
		log.Fatal(err)
	}

	scraperName, err := c.Flags().GetString("scraper")
	if err != nil {
		log.Fatal(err)
	}

	e.initConfig()

	var s scraper.IScraper
	conf := e.cfg.GetScraper(scraperName)
	switch scraperName {
	case "jav":
		s = scraper.NewJavLibraryScraper(proxy)

	case scraper.JavDB:
		s = scraper.NewJavDBScraper(conf)

	default:
		log.Fatalf("unknown scraper %s", scraperName)
	}

	for _, arg := range args {
		if err := s.Fetch(arg); err != nil {
			log.Fatal(err)
		}

		fmt.Println("number: ", s.GetNumber())
		fmt.Println("actors: ", s.GetActors())
		fmt.Println("cover: ", s.GetCover())
		fmt.Println("title: ", s.GetTitle())
		fmt.Println("studio: ", s.GetStudio())
		fmt.Println("release: ", s.GetRelease())
	}
}
