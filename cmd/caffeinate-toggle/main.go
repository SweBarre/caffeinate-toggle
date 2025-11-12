package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SweBarre/caffeinate-toggle/internal/tray"
	"github.com/SweBarre/caffeinate-toggle/internal/caffeinate"
)


var Version = "dev"

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println("Received signal:", sig)
		// Place your cleanup logic here!
		caffeinate.Stop()
		os.Exit(0)
	}()

	log.Printf("Starting CaffeniateToggle: %v", Version)
	tray.Run()
}
