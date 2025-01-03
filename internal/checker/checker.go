package checker

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hightemp/proxy_parser_checker/internal/config"
	"github.com/hightemp/proxy_parser_checker/internal/logger"
	"github.com/hightemp/proxy_parser_checker/internal/models/proxy"
)

var (
	mtx          sync.Mutex
	CheckRate    float32 = 0
	checkCounter int     = 0
)

type ProxyChecker struct {
	proxyChan  chan *proxy.Proxy
	maxWorkers int
}

func NewProxyChecker(cfg *config.Config) *ProxyChecker {
	maxWorkers := cfg.CheckerMaxWorkers
	if maxWorkers == 0 {
		maxWorkers = runtime.NumCPU() * 4
	}

	pc := &ProxyChecker{
		proxyChan:  make(chan *proxy.Proxy, maxWorkers*2),
		maxWorkers: maxWorkers,
	}

	for i := 0; i < maxWorkers; i++ {
		go pc.worker()
	}

	return pc
}

func (pc *ProxyChecker) worker() {
	for p := range pc.proxyChan {
		checkProxy(p)
	}
}

var checkURLs = []string{
	"https://api.ipify.org?format=json",
	"https://ifconfig.me/ip",
	"https://api.myip.com",
	"https://checkip.amazonaws.com",
}

type proxyCheckResult struct {
	success    bool
	pingTime   time.Duration
	detectedIP string
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

	var results []proxyCheckResult
	var totalPingTime time.Duration

	for _, checkURL := range checkURLs {
		result := checkSingleURL(client, checkURL)
		if result.success && result.detectedIP == lastProxy.Ip {
			results = append(results, result)
			totalPingTime += result.pingTime
		}
	}

	successRate := float64(len(results)) / float64(len(checkURLs))

	mtx.Lock()
	checkCounter++
	if successRate > 0.5 {
		lastProxy.IsWork = true
		lastProxy.PingTime = totalPingTime / time.Duration(len(results))
		lastProxy.SuccessCount++
		logger.LogInfo("[checker] Proxy checked successfully. Success rate: %.2f, Average ping time: %v",
			successRate, lastProxy.PingTime)
	} else {
		lastProxy.FailsCount++
		logger.LogError("[checker] Proxy check failed. Success rate: %.2f", successRate)
	}
	mtx.Unlock()

	if lastProxy.IsWork {
		logger.LogDebug("[checker][!] Found proxy: %s '%s:%s'", lastProxy.Protocol, lastProxy.Ip, lastProxy.Port)
		proxy.SaveWorkProxies()
		proxy.Save()
	}
}

func checkSingleURL(client *http.Client, checkURL string) proxyCheckResult {
	result := proxyCheckResult{success: false}

	startTime := time.Now()
	resp, err := client.Get(checkURL)
	result.pingTime = time.Since(startTime)

	if err != nil {
		logger.LogError("[checker] Request to %s failed: %v", checkURL, err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.LogError("[checker] Bad response status from %s: %d", checkURL, resp.StatusCode)
		return result
	}

	bodyBuffer := new(bytes.Buffer)
	_, err = bodyBuffer.ReadFrom(resp.Body)
	if err != nil {
		logger.LogError("[checker] Can't read body from %s: %v", checkURL, err)
		return result
	}

	if ip := extractIP(bodyBuffer.String(), checkURL); ip != "" {
		result.success = true
		result.detectedIP = ip
		logger.LogInfo("[checker] Successfully checked %s, detected IP: %s", checkURL, ip)
	}

	return result
}

func extractIP(body, checkURL string) string {
	switch {
	case checkURL == "https://api.ipify.org?format=json":
		var response map[string]interface{}
		if err := json.Unmarshal([]byte(body), &response); err == nil {
			if ip, ok := response["ip"].(string); ok {
				return ip
			}
		}
	case checkURL == "https://api.myip.com":
		var response map[string]interface{}
		if err := json.Unmarshal([]byte(body), &response); err == nil {
			if ip, ok := response["ip"].(string); ok {
				return ip
			}
		}
	default:
		return strings.TrimSpace(body)
	}
	return ""
}

func Loop(cfg *config.Config) {
	go func() {
		t := time.NewTicker(60 * time.Second)
		for range t.C {
			mtx.Lock()
			CheckRate = float32(checkCounter) / 60
			checkCounter = 0
			logger.LogDebug("[checker] Check rate: %f proxies per second", CheckRate)
			mtx.Unlock()
		}
	}()
	proxy.SetCheckPeriodDuration(cfg.CheckPeriodDuration)

	pc := NewProxyChecker(cfg)

	for {
		mtx.Lock()
		lastProxy := proxy.GetLastNotCheckedOne()
		mtx.Unlock()

		if lastProxy == nil {
			logger.LogDebug("[checker] No proxy found")
			time.Sleep(10 * time.Second)
			continue
		}

		logger.LogDebug("[checker] Checking proxy: %s '%s:%s'", lastProxy.Protocol, lastProxy.Ip, lastProxy.Port)
		pc.proxyChan <- lastProxy
	}
}
