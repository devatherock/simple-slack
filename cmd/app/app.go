//go:build !plugin && !integration
// +build !plugin,!integration

package main

import (
	"net/http"
	"os"

	"github.com/devatherock/vela-template-tester/pkg/util"
	log "github.com/sirupsen/logrus"
)

// Initializes log level
func init() {
	util.InitLogLevel()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	http.HandleFunc("/api/notification", sendNotification)
	http.HandleFunc("/api/health", checkHealth)

	http.ListenAndServe(":"+lookupPort(), nil)
}

// Reads port from PORT environment variable
func lookupPort() string {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	log.Info("Listening on port ", port)

	return port
}
