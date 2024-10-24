package proxy

import (
	"fmt"
	"os"
	"time"

	"github.com/hightemp/proxy_parser_checker/internal/logger"
	"gopkg.in/yaml.v3"
)

const (
	PROTO_HTTP  = "http"
	PROTO_HTTPS = "https"
)

type Proxy struct {
	Ip              string    `yaml:"ip"`
	Port            string    `yaml:"port"`
	Protocol        string    `yaml:"protocol"`
	LastCheckedTime time.Time `yaml:"last_checked_time"`
}

var (
	IsDirty     bool = false
	proxiesList []Proxy
)

func Find(p Proxy) int {
	for i, pi := range proxiesList {
		if pi.Ip == p.Ip && pi.Port == p.Port {
			return i
		}
	}

	return -1
}

func Add(p Proxy) {
	index := Find(p)

	if index == -1 {
		proxiesList = append(proxiesList, p)
		IsDirty = true
		logger.LogDebug("[proxy] added proxy '%s:%s'", p.Ip, p.Port)
	}
}

func AddList(pl []Proxy) {
	for _, p := range pl {
		Add(p)
	}
}

func Save() error {
	if !IsDirty {
		return nil
	}

	yamlText, err := yaml.Marshal(proxiesList)

	if err != nil {
		return fmt.Errorf("Can't pack to yaml: %v", err)
	}

	err = os.WriteFile("./out/all_proxies.yaml", yamlText, 0644)

	if err != nil {
		return fmt.Errorf("Can't write file: %v", err)
	}

	IsDirty = false
	return nil
}
