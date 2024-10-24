package parsers

import (
	"encoding/json"
	"strings"

	"github.com/hightemp/proxy_parser_checker/internal/logger"
	"github.com/hightemp/proxy_parser_checker/internal/models/proxy"
)

/*
{
  "data": [
    {
      "_id": "66052aae6fb9cbee3795eacc",
      "ip": "124.6.225.124",
      "anonymityLevel": "elite",
      "asn": "AS38256",
      "city": "Dhaka",
      "country": "BD",
      "created_at": "2024-03-28T08:30:38.985Z",
      "google": false,
      "isp": "Bengal Group",
      "lastChecked": 1729670444,
      "latency": 162.818,
      "org": "Prisma Digital Network Ltd.",
      "port": "1088",
      "protocols": [
        "socks4"
      ],
      "speed": 2,
      "upTime": 95.43348775645268,
      "upTimeSuccessCount": 1442,
      "upTimeTryCount": 1511,
      "updated_at": "2024-10-23T08:00:44.060Z",
      "responseTime": 2396
    },
*/

type ProxyListParser struct{}

func (p *ProxyListParser) IsTargetSite(url string) bool {
	return strings.Contains(url, "proxylist.geonode.com/api/proxy-list")
}

func (p *ProxyListParser) ParseProxyList(s string) []proxy.Proxy {
	var proxyList []proxy.Proxy
	var o map[string]interface{}

	err := json.Unmarshal([]byte(s), &o)

	if err != nil {
		logger.LogError("Can't parse json: %v", err)
		return proxyList
	}

	data, ok := o["data"].([]interface{})
	if !ok {
		logger.LogError("No 'data' field found or invalid format: %v", err)
		return proxyList
	}

	for _, item := range data {
		proxyMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		ip, ok := proxyMap["ip"].(string)
		if !ok {
			continue
		}

		port, ok := proxyMap["port"].(string)
		if !ok {
			continue
		}

		protocols, ok := proxyMap["protocols"].([]interface{})
		if !ok {
			continue
		}

		hasValidProtocol := false
		proxyType := ""

		for _, proto := range protocols {
			if protoStr, ok := proto.(string); ok {
				protoLower := strings.ToLower(protoStr)
				if protoLower == "http" || protoLower == "https" {
					hasValidProtocol = true
					proxyType = protoLower
					break
				}
			}
		}

		if !hasValidProtocol {
			continue
		}

		proxyList = append(proxyList, proxy.Proxy{Ip: ip, Port: port, Protocol: proxyType})
	}

	return proxyList
}
