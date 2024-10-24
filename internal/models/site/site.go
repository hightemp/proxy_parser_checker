package site

import (
	"time"

	"github.com/hightemp/proxy_parser_checker/internal/config"
)

type Site struct {
	Url            string
	LastParsedTime time.Time
}

var cfg *config.Config
var sites []Site

func Init(c *config.Config) {
	cfg = c
}

func AddSites() {
	for _, url := range cfg.SitesForParsing {
		sites = append(sites, Site{Url: url})
	}
}

func IsExpired(t time.Time) bool {
	now := time.Now()
	expirationTime := t.Add(cfg.ParsePeriodDuration)
	return now.After(expirationTime)
}

func GetLastOne() *Site {
	for _, s := range sites {
		if IsExpired(s.LastParsedTime) {
			return &s
		}
	}

	return nil
}
