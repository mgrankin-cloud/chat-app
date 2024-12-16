package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mgrankin-cloud/messenger/cmd/load-balancer/utils"
	"github.com/mgrankin-cloud/messenger/internal/config/auth"
	"github.com/mgrankin-cloud/messenger/internal/config/chat"
	"github.com/mgrankin-cloud/messenger/internal/config/message"
)

func main() {
	authCfg := auth.MustLoad()
	chatCfg := chat.MustLoad()
	messageCfg := message.MustLoad()

	startLoadBalancer(authCfg, "auth")
	startLoadBalancer(chatCfg, "chat")
	startLoadBalancer(messageCfg, "message")

	// Keep the main function running
	select {}
}

func startLoadBalancer(cfg interface {
	ServerList          []string
	LBAddress           string
	HealthCheckInterval time.Duration
}, serviceName string) {
	utils.ServerList = cfg.ServerList
	utils.LBAddress = cfg.LBAddress

	// Start health check routine
	go func() {
		ticker := time.NewTicker(cfg.HealthCheckInterval)
		for {
			utils.CheckHealth()
			<-ticker.C
		}
	}()

	// Start load balancer server
	lbServer := &http.Server{
		Addr:    utils.LBAddress,
		Handler: http.HandlerFunc(utils.LBHandler),
	}

	log.Printf("Starting load balancer for %s service on %s", serviceName, utils.LBAddress)
	go func() {
		log.Fatal(lbServer.ListenAndServe())
	}()
}