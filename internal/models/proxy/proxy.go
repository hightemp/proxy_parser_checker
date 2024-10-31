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
	Ip              string        `yaml:"ip"`
	Port            string        `yaml:"port"`
	Protocol        string        `yaml:"protocol"`
	LastCheckedTime time.Time     `yaml:"last_checked_time"`
	PingTime        time.Duration `yaml:"ping_time"`
	IsWork          bool          `yaml:"is_work"`
}

var (
	IsDirty             bool = false
	proxiesList         []Proxy
	checkPeriodDuration time.Duration
)

func SetCheckPeriodDuration(t time.Duration) {
	checkPeriodDuration = t
}

func Find(p Proxy) int {
	for i, pi := range proxiesList {
		if pi.Ip == p.Ip && pi.Port == p.Port {
			return i
		}
	}

	return -1
}

func Delete(p Proxy) bool {
	index := Find(p)
	if index != -1 {
		proxiesList = append(proxiesList[:index], proxiesList[index+1:]...)
		IsDirty = true
		logger.LogDebug("[proxy] deleted proxy '%s:%s'", p.Ip, p.Port)
		return true
	}
	return false
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

func IsExpired(t time.Time) bool {
	now := time.Now()
	expirationTime := t.Add(checkPeriodDuration)
	return now.After(expirationTime)
}

func GetLastNotCheckedOne() *Proxy {
	for i := range proxiesList {
		if IsExpired(proxiesList[i].LastCheckedTime) {
			return &proxiesList[i]
		}
	}

	return nil
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

func GetWorkProxies() []*Proxy {
	var result []*Proxy

	for i := range proxiesList {
		if proxiesList[i].IsWork {
			result = append(result, &proxiesList[i])
		}
	}

	return result
}

func SaveWorkProxies() error {
	yamlText, err := yaml.Marshal(GetWorkProxies())

	if err != nil {
		return fmt.Errorf("Can't pack to yaml: %v", err)
	}

	err = os.WriteFile("./out/work_proxies.yaml", yamlText, 0644)

	if err != nil {
		return fmt.Errorf("Can't write file: %v", err)
	}

	IsDirty = false
	return nil
}

func GetAllProxies() []Proxy {
	return proxiesList
}
