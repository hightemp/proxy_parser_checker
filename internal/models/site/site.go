package site

import (
	"fmt"
	"os"
	"time"

	"github.com/hightemp/proxy_parser_checker/internal/logger"
	"gopkg.in/yaml.v3"
)

type Site struct {
	Url            string
	LastParsedTime time.Time
}

var (
	sites               []Site
	IsDirty             = false
	parsePeriodDuration time.Duration
)

func SetParsePeriodDuration(t time.Duration) {
	parsePeriodDuration = t
}

func FindUrl(url string) int {
	for i, si := range sites {
		if si.Url == url {
			return i
		}
	}

	return -1
}

func Add(url string) {
	index := FindUrl(url)

	if index == -1 {
		sites = append(sites, Site{Url: url})
		IsDirty = true
		logger.LogDebug("[site] added site '%s'", url)
	}
}

func AddList(urlList []string) {
	for _, url := range urlList {
		Add(url)
	}
}

func IsExpired(t time.Time) bool {
	now := time.Now()
	expirationTime := t.Add(parsePeriodDuration)
	return now.After(expirationTime)
}

func GetLastOne() *Site {
	for i := range sites {
		if IsExpired(sites[i].LastParsedTime) {
			return &sites[i]
		}
	}

	return nil
}

func Save() error {
	if !IsDirty {
		return nil
	}

	yamlText, err := yaml.Marshal(sites)

	if err != nil {
		return fmt.Errorf("Can't pack to yaml: %v", err)
	}

	err = os.WriteFile("./sites_for_parsing.yaml", yamlText, 0644)

	if err != nil {
		return fmt.Errorf("Can't write file: %v", err)
	}

	IsDirty = false
	return nil
}
