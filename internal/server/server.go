package server

import (
	"encoding/json"
	"net/http"

	"github.com/hightemp/proxy_parser_checker/internal/checker"
	"github.com/hightemp/proxy_parser_checker/internal/config"
	"github.com/hightemp/proxy_parser_checker/internal/logger"
	"github.com/hightemp/proxy_parser_checker/internal/models/proxy"
	"github.com/hightemp/proxy_parser_checker/internal/models/site"
)

type ProxyResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

func jsonResponse(w http.ResponseWriter, status int, resp ProxyResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func handleProxies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		jsonResponse(w, http.StatusOK, ProxyResponse{
			Success: true,
			Data:    proxy.GetAllProxies(),
		})
	case http.MethodPost:
		var newProxy proxy.Proxy
		if err := json.NewDecoder(r.Body).Decode(&newProxy); err != nil {
			jsonResponse(w, http.StatusBadRequest, ProxyResponse{
				Success: false,
				Error:   "Invalid request body",
			})
			return
		}

		proxy.Add(newProxy)
		if err := proxy.Save(); err != nil {
			jsonResponse(w, http.StatusInternalServerError, ProxyResponse{
				Success: false,
				Error:   "Failed to save proxy",
			})
			return
		}

		jsonResponse(w, http.StatusCreated, ProxyResponse{
			Success: true,
			Data:    newProxy,
		})
	case http.MethodDelete:
		var proxyToDelete proxy.Proxy
		if err := json.NewDecoder(r.Body).Decode(&proxyToDelete); err != nil {
			jsonResponse(w, http.StatusBadRequest, ProxyResponse{
				Success: false,
				Error:   "Invalid request body",
			})
			return
		}

		if deleted := proxy.Delete(proxyToDelete); !deleted {
			jsonResponse(w, http.StatusNotFound, ProxyResponse{
				Success: false,
				Error:   "Proxy not found",
			})
			return
		}

		if err := proxy.Save(); err != nil {
			jsonResponse(w, http.StatusInternalServerError, ProxyResponse{
				Success: false,
				Error:   "Failed to save changes",
			})
			return
		}

		jsonResponse(w, http.StatusOK, ProxyResponse{
			Success: true,
			Data:    "Proxy deleted successfully",
		})
	default:
		jsonResponse(w, http.StatusMethodNotAllowed, ProxyResponse{
			Success: false,
			Error:   "Method not allowed",
		})
	}
}

func handleSites(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		jsonResponse(w, http.StatusOK, ProxyResponse{
			Success: true,
			Data:    site.GetAllSites(),
		})
	case http.MethodPost:
		var requestBody struct {
			URL string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			jsonResponse(w, http.StatusBadRequest, ProxyResponse{
				Success: false,
				Error:   "Invalid request body",
			})
			return
		}

		if requestBody.URL == "" {
			jsonResponse(w, http.StatusBadRequest, ProxyResponse{
				Success: false,
				Error:   "URL is required",
			})
			return
		}

		site.Add(requestBody.URL)
		if err := site.Save(); err != nil {
			jsonResponse(w, http.StatusInternalServerError, ProxyResponse{
				Success: false,
				Error:   "Failed to save site",
			})
			return
		}

		jsonResponse(w, http.StatusCreated, ProxyResponse{
			Success: true,
			Data:    requestBody.URL,
		})
	case http.MethodDelete:
		var requestBody struct {
			URL string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			jsonResponse(w, http.StatusBadRequest, ProxyResponse{
				Success: false,
				Error:   "Invalid request body",
			})
			return
		}

		if requestBody.URL == "" {
			jsonResponse(w, http.StatusBadRequest, ProxyResponse{
				Success: false,
				Error:   "URL is required",
			})
			return
		}

		if deleted := site.Delete(requestBody.URL); !deleted {
			jsonResponse(w, http.StatusNotFound, ProxyResponse{
				Success: false,
				Error:   "Site not found",
			})
			return
		}

		if err := site.Save(); err != nil {
			jsonResponse(w, http.StatusInternalServerError, ProxyResponse{
				Success: false,
				Error:   "Failed to save changes",
			})
			return
		}

		jsonResponse(w, http.StatusOK, ProxyResponse{
			Success: true,
			Data:    "Site deleted successfully",
		})
	default:
		jsonResponse(w, http.StatusMethodNotAllowed, ProxyResponse{
			Success: false,
			Error:   "Method not allowed",
		})
	}
}

func handleWorkingProxies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, ProxyResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}

	jsonResponse(w, http.StatusOK, ProxyResponse{
		Success: true,
		Data:    proxy.GetWorkProxies(),
	})
}

func handleFirstWorkingProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, ProxyResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}

	workProxies := proxy.GetWorkProxies()
	if len(workProxies) == 0 {
		jsonResponse(w, http.StatusNotFound, ProxyResponse{
			Success: false,
			Error:   "No working proxies available",
		})
		return
	}

	jsonResponse(w, http.StatusOK, ProxyResponse{
		Success: true,
		Data:    workProxies[0],
	})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, ProxyResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}

	type StatsInfo struct {
		TotalProxies      int     `json:"total_proxies"`
		WorkedProxies     int     `json:"worked_proxies"`
		BlockedProxies    int     `json:"blocked_proxies"`
		NotCheckedProxies int     `json:"not_checked_proxies"`
		CheckRate         float32 `json:"check_rate"`
	}

	proxies := proxy.GetAllProxies()
	workedProxies, blockedProxies, notCheckedProxies := 0, 0, 0

	for _, p := range proxies {
		if p.IsWork && p.FailsCount < proxy.MaxFailsCount {
			workedProxies++
		}
		if p.FailsCount >= proxy.MaxFailsCount {
			blockedProxies++
		}
		if p.FailsCount < proxy.MaxFailsCount && proxy.IsExpired(p.LastCheckedTime) {
			notCheckedProxies++
		}
	}

	statsInfo := StatsInfo{
		TotalProxies:      len(proxies),
		WorkedProxies:     workedProxies,
		BlockedProxies:    blockedProxies,
		NotCheckedProxies: notCheckedProxies,
		CheckRate:         checker.CheckRate,
	}

	jsonResponse(w, http.StatusOK, ProxyResponse{
		Success: true,
		Data:    statsInfo,
	})
}

func Start() {
	http.HandleFunc("/api/v1/proxies", handleProxies)
	http.HandleFunc("/api/v1/proxies/working", handleWorkingProxies)
	http.HandleFunc("/api/v1/proxies/working/first", handleFirstWorkingProxy)

	http.HandleFunc("/api/v1/sites", handleSites)

	http.HandleFunc("/api/v1/stats", handleStats)

	port := config.GetConfig().ServerPort
	if port == "" {
		port = "8080"
	}

	logger.LogInfo("Starting HTTP server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.PanicError("Failed to start HTTP server: %v", err)
	}
}
