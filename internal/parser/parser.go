package parser

import (
	"bytes"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/hightemp/proxy_parser_checker/internal/config"
	"github.com/hightemp/proxy_parser_checker/internal/logger"
	"github.com/hightemp/proxy_parser_checker/internal/models/proxy"
	"github.com/hightemp/proxy_parser_checker/internal/models/site"
	"github.com/hightemp/proxy_parser_checker/internal/parser/parsers"
)

type IParser interface {
	IsTargetSite(url string) bool
	ParseProxyList(s string) []proxy.Proxy
}

var (
	parsersList []IParser
	mtx         sync.Mutex
)

func AddParser(p IParser) {
	parsersList = append(parsersList, p)
}

func parseSite(client *http.Client, lastSite *site.Site) {
	mtx.Lock()
	defer mtx.Unlock()

	logger.LogDebug("[parser] Making request to '%s'", lastSite.Url)
	resp, err := client.Get(lastSite.Url)

	if err != nil {
		logger.LogError("[parser] Can't get url: '%s', %v", lastSite.Url, err)
		lastSite.LastParsedTime = time.Now()
		return
	}

	bodyBuffer := new(bytes.Buffer)
	_, err = bodyBuffer.ReadFrom(resp.Body)
	if err != nil {
		logger.LogError("[parser] Can't read body: %v", err)
		lastSite.LastParsedTime = time.Now()
		return
	}
	body := bodyBuffer.String()
	resp.Body.Close()

	// fmt.Println(body)

	logger.LogDebug("[parser] parsing '%s'", lastSite.Url)

	for _, p := range parsersList {
		if p.IsTargetSite(lastSite.Url) {
			logger.LogDebug("[parser] detected '%s'", reflect.TypeOf(p).String())
			proxy.AddList(p.ParseProxyList(body))
			lastSite.LastParsedTime = time.Now()
		}
	}
}

func ParsingLoop(c *config.Config) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	AddParser(&parsers.ProxyListParser{})
	AddParser(&parsers.TextListParser{})

	site.Init(c)
	site.AddSites()

	for {
		lastSite := site.GetLastOne()

		if lastSite == nil {
			time.Sleep(time.Minute)
			proxy.Save()
			continue
		}

		go parseSite(client, lastSite)
	}
}
