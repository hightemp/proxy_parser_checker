package server

import (
	"encoding/json"
	"net/http"

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

func handleWorkedProxies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, ProxyResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}

	workProxies := proxy.GetWorkProxies()
	jsonResponse(w, http.StatusOK, ProxyResponse{
		Success: true,
		Data:    workProxies,
	})
}

func handleGetWorkProxy(w http.ResponseWriter, r *http.Request) {
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

	// Return the first working proxy
	jsonResponse(w, http.StatusOK, ProxyResponse{
		Success: true,
		Data:    workProxies[0],
	})
}

func handleGetAllProxies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, ProxyResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}

	jsonResponse(w, http.StatusOK, ProxyResponse{
		Success: true,
		Data:    proxy.GetAllProxies(),
	})
}

func handleAddProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, ProxyResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}

	var newProxy proxy.Proxy
	if err := json.NewDecoder(r.Body).Decode(&newProxy); err != nil {
		jsonResponse(w, http.StatusBadRequest, ProxyResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	proxy.Add(newProxy)
	proxy.Save()
	if err := proxy.Save(); err != nil {
		jsonResponse(w, http.StatusInternalServerError, ProxyResponse{
			Success: false,
			Error:   "Failed to save proxy",
		})
		return
	}

	jsonResponse(w, http.StatusOK, ProxyResponse{
		Success: true,
		Data:    newProxy,
	})
}

func handleGetAllSites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, ProxyResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}

	jsonResponse(w, http.StatusOK, ProxyResponse{
		Success: true,
		Data:    site.GetAllSites(),
	})
}

func handleAddSite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, ProxyResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}

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
	site.Save()
	if err := site.Save(); err != nil {
		jsonResponse(w, http.StatusInternalServerError, ProxyResponse{
			Success: false,
			Error:   "Failed to save site",
		})
		return
	}

	jsonResponse(w, http.StatusOK, ProxyResponse{
		Success: true,
		Data:    requestBody.URL,
	})
}

func Start() {
	http.HandleFunc("/work-proxies/all", handleWorkedProxies)
	http.HandleFunc("/work-proxies/one", handleGetWorkProxy)
	http.HandleFunc("/proxies/all", handleGetAllProxies)
	http.HandleFunc("/proxies/add", handleAddProxy)
	http.HandleFunc("/sites/all", handleGetAllSites)
	http.HandleFunc("/sites/add", handleAddSite)

	port := config.GetConfig().ServerPort
	if port == "" {
		port = "8080"
	}

	logger.LogInfo("Starting HTTP server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.PanicError("Failed to start HTTP server: %v", err)
	}
}
