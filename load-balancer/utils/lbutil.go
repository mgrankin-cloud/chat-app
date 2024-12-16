package utils

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func CheckHealth() {
	t := time.Now()
	var activeServers []string

	for _, serverURL := range serverList {
		resp, err := http.Get(serverURL + "/health")
		if err != nil {
			log.Printf("Server %s is not healthy: %v", serverURL, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Server %s is not healthy: status code %d", serverURL, resp.StatusCode)
			continue
		}

		activeServers = append(activeServers, serverURL)
	}

	mu.Lock()
	serverList = activeServers
	mu.Unlock()

	log.Printf("Health check performed @ %d:%d:%d", t.Hour(), t.Minute(), t.Second())
}

func GetTargetURL() *url.URL {
	mu.Lock()
	defer mu.Unlock()

	if len(serverList) == 0 {
		log.Fatal("No available servers")
	}

	serverIndex := selectedServer % len(serverList)
	target, err := url.Parse(serverList[serverIndex])
	if err != nil {
		log.Fatalf("Error parsing URL: %s", err.Error())
	}
	selectedServer++
	return target
}

func LBHandler(w http.ResponseWriter, req *http.Request) {
	target := GetTargetURL()
	log.Printf("Routing request to %s: %s %s", target, req.Method, req.URL.Path)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, req)
}