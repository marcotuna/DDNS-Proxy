package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"ddns-proxy/pkg/cloudflare"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "0.0.0.0", 80),
		Handler: r,
	}

	log.Infof("Webserver listening on %s", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Listen: %s", err)
			return
		}
	}()

	r.GET("/update", func(c *gin.Context) {
		api, err := cloudflare.Setup(os.Getenv("API_KEY"), os.Getenv("API_EMAIL"))

		if err != nil {
			log.Fatal(err)
			return
		}

		// Initialize Method Struct
		cf := cloudflare.Hub{API: api}

		// Fetch Zone ID
		zoneID, err := cf.ZoneID("sampledomain.com")

		if err != nil {
			log.Fatal(err)
			return
		}

		// Create or Update DNS
		createOrUpdateDNSRecord, err := cf.CreateOrUpdateDNSRecord(zoneID, "A", "sampledomain", "127.0.0.8")

		if err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(200, gin.H{
			"status": createOrUpdateDNSRecord,
		})
	})

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Trace("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server Shutdown: %s", err)
	}

	log.Trace("Server exiting")
}
