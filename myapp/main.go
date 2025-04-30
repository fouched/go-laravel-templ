package main

import (
	"context"
	"fmt"
	"github.com/fouched/rapidus"
	"myapp/data"
	"myapp/handlers"
	"os"
	"os/signal"
	"syscall"
)

type application struct {
	App      *rapidus.Rapidus
	Handlers *handlers.Handlers
	Models   data.Models
}

func main() {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the application
	go func() {
		r := initApplication()
		r.App.ListenAndServe()
	}()

	// Listen for the interrupt signal to complete, giving defer statements time to run
	<-ctx.Done()
	fmt.Println("Graceful shutdown complete")
}
