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

func checkProxy(lastProxy *proxy.Proxy) {
	mtx.Lock()
	defer mtx.Unlock()

	lastProxy.LastCheckedTime = time.Now()
	lastProxy.IsWork = false

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
		Timeout:   time.Second * 30,
	}

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

	lastProxy.IsWork = true
	lastProxy.PingTime = pingTime

	logger.LogInfo("[checker] Proxy checked successfully. Ping time: %v", pingTime)
}

func Loop(cfg *config.Config) {
	proxy.SetCheckPeriodDuration(cfg.CheckPeriodDuration)

	for {
		lastProxy := proxy.GetLastNotCheckedOne()

		if lastProxy == nil {
			logger.LogDebug("[checker] No proxy found")
			time.Sleep(10 * time.Second)
			continue
		}

		logger.LogDebug("[checker] Found proxy: %s '%s:%s'", lastProxy.Protocol, lastProxy.Ip, lastProxy.Port)
		go checkProxy(lastProxy)
	}
}
