package main

import (
	"github.com/mgrankin-cloud/messenger/cmd/load-balancer/utils"

	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/serverSSO", sso.ssoServerHandler)
	r.HandleFunc("/serverUser", user.userServerHandler)
	r.HandleFunc("/serverChat", chat.chatServerHandler)
	r.HandleFunc("/serverMessage", message.messageServerHandler)
	r.HandleFunc("/serverModels", models.modelsServerHandler)
	r.HandleFunc("/serverMedia", media.mediaServerHandler)
	r.HandleFunc("/serverNotification", notification.notificationServerHandler)

	go func() {
		log.Fatal(server.Server1.ListenAndServe())
	}()

	go func() {
		log.Fatal(server.Server2.ListenAndServe())
	}()

	go func() {
		var ticker *time.Ticker = time.NewTicker(15 * time.Second)
		for {
			utils.CheckHealth()
			<-ticker.C
		}
	}()

	lbServer := &http.Server{
		Addr:    utils.LBAddress,
		Handler: http.HandlerFunc(utils.LBHandler),
	}

	log.Fatal(lbServer.ListenAndServe())
	select {}
}
