package main

import (
	"log"
	"mail/config"
	"mail/delivery/http"
	"mail/delivery/worker"
	"os"
)

func main() {

	envSource := "SYSTEM"

	if os.Getenv("BYPASS_ENV_FILE") == "" {
		log.Println("[INFO] Load Config")
		config.LoadConfig()
		envSource = "FILE"
	}

	log.Println("[INFO] Loaded Config : " + envSource)

	if os.Getenv("APPLICATION_MODE") == "worker" {
		worker.RunWorker()
	} else {
		http.RunAPI()
	}
}
