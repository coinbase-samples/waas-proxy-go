package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/coinbase-samples/waas-proxy-go/config"
	"github.com/coinbase-samples/waas-proxy-go/handlers"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {

	var app config.AppConfig

	if err := config.Setup(&app); err != nil {
		log.Fatalf("unable to config app: %v", err)
	}

	config.LogInit(app)

	wait := time.Minute * 1
	if len(os.Getenv("GRACEFUL_TIMEOUT")) > 0 {
		var err error
		if wait, err = time.ParseDuration(os.Getenv("GRACEFUL_TIMEOUT")); err != nil {
			log.Fatalf("Invalid GRACEFUL_TIMEOUT: %s - err: %v", os.Getenv("GRACEFUL_TIMEOUT"), err)
		}
	}

	if err := waas.InitClients(app); err != nil {
		log.Fatalf("unable to init WaaS clients: %v", err)
	}

	router := mux.NewRouter()

	handlers.RegisterHandlers(app, router)

	port := "8443"

	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}

	headersOk := ghandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := ghandlers.AllowedOrigins([]string{app.AppUrl})
	methodsOk := ghandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	srv := &http.Server{
		Handler:      ghandlers.CORS(originsOk, headersOk, methodsOk)(router),
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	log.Infof(fmt.Sprintf("Starting listener on: %s", port))

	go func() {
		if app.IsLocalEnv() {
			if err := srv.ListenAndServe(); err != nil {
				log.Fatalf("ListenAndServe: %v", err)
			}

		} else {
			if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != nil {
				log.Fatalf("ListenAndServeTLS: %v", err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)

}
