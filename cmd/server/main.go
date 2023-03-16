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
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {

	var app config.AppConfig

	if err := config.Setup(&app); err != nil {
		log.Fatalf("Unable to config app: %v", err)
	}

	config.LogInit(app)

	wait := time.Minute * 1
	if len(os.Getenv("GRACEFUL_TIMEOUT")) > 0 {
		var err error
		if wait, err = time.ParseDuration(os.Getenv("GRACEFUL_TIMEOUT")); err != nil {
			log.Fatalf("Invalid GRACEFUL_TIMEOUT: %s - err: %v", os.Getenv("GRACEFUL_TIMEOUT"), err)
		}
	}

	router := mux.NewRouter()

	handlers.RegisterHandlers(app, router)

	port := "8443"

	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}

	log.Infof(fmt.Sprintf("starting listener on: %s", port))

	headersOk := ghandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := ghandlers.AllowedOrigins([]string{"https://app.wenthemerge.xyz"})
	methodsOk := ghandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	srv := &http.Server{
		Handler:      ghandlers.CORS(originsOk, headersOk, methodsOk)(router),
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)

}
