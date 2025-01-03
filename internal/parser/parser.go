package parser

import (
	"bytes"
	"net/http"
	"reflect"
	"runtime"
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

type WorkerPool struct {
	siteChan   chan *site.Site
	wg         sync.WaitGroup
	client     *http.Client
	maxWorkers int
}

var (
	parsersList []IParser
	mtx         sync.Mutex
)

func AddParser(p IParser) {
	parsersList = append(parsersList, p)
}

func NewWorkerPool(cfg *config.Config) *WorkerPool {
	maxWorkers := cfg.ParserMaxWorkers
	if maxWorkers == 0 {
		maxWorkers = runtime.NumCPU() * 4
	}

	wp := &WorkerPool{
		siteChan: make(chan *site.Site, maxWorkers),
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		maxWorkers: maxWorkers,
	}

	wp.StartWorkers()

	return wp
}

func (w *WorkerPool) StartWorkers() {
	for i := 0; i < w.maxWorkers; i++ {
		w.wg.Add(1)
		go w.work()
	}
}

func (w *WorkerPool) work() {
	defer w.wg.Done()

	for s := range w.siteChan {
		w.parse(s)
	}
}

func (w *WorkerPool) parse(lastSite *site.Site) {
	mtx.Lock()
	lastSite.LastParsedTime = time.Now()
	mtx.Unlock()

	logger.LogDebug("[parser] Making request to '%s'", lastSite.Url)
	resp, err := w.client.Get(lastSite.Url)

	if err != nil {
		logger.LogError("[parser] Can't get url: '%s', %v", lastSite.Url, err)
		return
	}

	bodyBuffer := new(bytes.Buffer)
	_, err = bodyBuffer.ReadFrom(resp.Body)
	if err != nil {
		logger.LogError("[parser] Can't read body: %v", err)
		return
	}
	body := bodyBuffer.String()
	resp.Body.Close()

	logger.LogDebug("[parser] parsing '%s'", lastSite.Url)

	for _, p := range parsersList {
		if p.IsTargetSite(lastSite.Url) {
			logger.LogDebug("[parser] detected '%s'", reflect.TypeOf(p).String())
			mtx.Lock()
			proxy.AddList(p.ParseProxyList(body))
			proxy.Save()
			mtx.Unlock()
			return
		}
	}
}

func Loop(cfg *config.Config) {
	w := NewWorkerPool(cfg)

	AddParser(&parsers.ProxyListParser{})
	AddParser(&parsers.TextListParser{})

	site.SetParsePeriodDuration(cfg.ParsePeriodDuration)
	if site.FileExists() {
		site.Load()
	} else {
		site.AddList(cfg.SitesForParsing)
		site.Save()
	}

	for {
		lastSite := site.GetLastOne()

		if lastSite == nil {
			logger.LogDebug("[parser] No site found")
			time.Sleep(time.Minute)
			continue
		}

		logger.LogDebug("[parser] Found site: '%s'", lastSite.Url)
		w.siteChan <- lastSite
	}
}
