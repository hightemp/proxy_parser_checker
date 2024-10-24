package parsers

import (
	"regexp"
	"strings"

	"github.com/hightemp/proxy_parser_checker/internal/models/proxy"
)

type TextListParser struct{}

func (p *TextListParser) IsTargetSite(url string) bool {
	return true
}

func (p *TextListParser) ParseProxyList(s string) []proxy.Proxy {
	var proxyList []proxy.Proxy

	ipPortRegex := regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d{1,5})$`)

	lines := strings.Split(strings.TrimSpace(s), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		matches := ipPortRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		ip := matches[1]
		port := matches[2]

		proxyList = append(proxyList, proxy.Proxy{Ip: ip, Port: port, Protocol: proxy.PROTO_HTTP})
	}

	return proxyList
}
