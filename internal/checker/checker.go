package checker

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/hightemp/proxy_parser_checker/internal/config"
	"github.com/hightemp/proxy_parser_checker/internal/logger"
	"github.com/hightemp/proxy_parser_checker/internal/models/proxy"
)

var (
	mtx sync.Mutex
)

const (
	maxWorkers = 10
)

type ProxyChecker struct {
	proxyChan chan *proxy.Proxy
	wg        sync.WaitGroup
}

func NewProxyChecker() *ProxyChecker {
	return &ProxyChecker{
		proxyChan: make(chan *proxy.Proxy, maxWorkers),
	}
}

func (pc *ProxyChecker) worker() {
	defer pc.wg.Done()

	for p := range pc.proxyChan {
		checkProxy(p)
	}
}

func checkProxy(lastProxy *proxy.Proxy) {
	mtx.Lock()
	lastProxy.LastCheckedTime = time.Now()
	lastProxy.IsWork = false
	mtx.Unlock()

	proxyURL := fmt.Sprintf("%s://%s:%s", lastProxy.Protocol, lastProxy.Ip, lastProxy.Port)
	proxyUrlParsed, err := url.Parse(proxyURL)
	if err != nil {
		logger.LogError("[checker] Failed to parse proxy URL: %v", err)
		return
	}

	transport := &http.Transport{
		Proxy:             http.ProxyURL(proxyUrlParsed),
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 5,
	}

	logger.LogDebug("[checker] Making request with proxy '%s://%s:%s'", lastProxy.Protocol, lastProxy.Ip, lastProxy.Port)
	startTime := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	pingTime := time.Since(startTime)

	if err != nil {
		logger.LogError("[checker] Proxy check failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.LogError("[checker] Bad response status: %d", resp.StatusCode)
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

	// {"ip":"5.1.1.1"}
	var o map[string]interface{}

	err = json.Unmarshal([]byte(body), &o)

	if err != nil {
		logger.LogError("Can't parse json: %v", err)
		return
	}

	ip, ok := o["ip"].(string)
	if !ok {
		logger.LogError("Can't parse ip: %v", err)
		return
	}

	logger.LogInfo("[checker] Found response ip: %v", ip)

	mtx.Lock()
	lastProxy.IsWork = true
	lastProxy.PingTime = pingTime
	mtx.Unlock()

	logger.LogInfo("[checker] Proxy checked successfully. Ping time: %v", pingTime)
	proxy.SaveWorkProxies()
	proxy.Save()
}

func Loop(cfg *config.Config) {
	proxy.SetCheckPeriodDuration(cfg.CheckPeriodDuration)

	pc := NewProxyChecker()

	for i := 0; i < maxWorkers; i++ {
		pc.wg.Add(1)
		go pc.worker()
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		lastProxy := proxy.GetLastNotCheckedOne()
		if lastProxy == nil {
			logger.LogDebug("[checker] No proxy found")
			time.Sleep(10 * time.Second)
			continue
		}

		logger.LogDebug("[checker] Found proxy: %s '%s:%s'", lastProxy.Protocol, lastProxy.Ip, lastProxy.Port)
		pc.proxyChan <- lastProxy
	}
}
