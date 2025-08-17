package main

import (
	"fmt"
	internalApp "test-server/internal/app"
)

func main() {
	// app, err := internalApp.NewApp(os.Getenv("CONFIG_FILE"))
	app, err := internalApp.NewApp("configs/config.yaml")
	if err != nil {
		fmt.Printf("main: error occured while starting the app: %v", err)
		return
	}

	if err := app.ListenAndServe(); err != nil {
		fmt.Printf("main: error occured while starting server: %v", err)
		return
	}
}
