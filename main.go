package main

import (
	"context"
	"embed"
	"log"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS
var App *application.App

var (
	healthMu     sync.Mutex
	healthCancel context.CancelFunc
)

func init() {
	// Register a custom event whose associated data type is string.
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	application.RegisterEvent[string]("errMsg")
	application.RegisterEvent[Scoop]("respMsg")
	application.RegisterEvent[Server]("initiateHealthCheck")
	application.RegisterEvent[string]("serverHealth")
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.

	SyncService := &SyncServer{}

	App = application.New(application.Options{
		Name:        "Scoop",
		Description: "REST API client for testing and discovery",
		Services: []application.Service{
			application.NewService(&ScoopService{}),
			application.NewService(SyncService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	App.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Scoop",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(10, 10, 10),
		URL:              "/",
	})

	// will cancel running go routine if called multiple times, only one health check runs at a time
	App.Event.On("initiateHealthCheck", func(e *application.CustomEvent) {
		server, ok := e.Data.(Server)
		if !ok {
			App.Event.Emit("errMsg", "initiateHealthCheck payload was not of type Server")
			App.Event.Emit("serverHealth", "Offline")
			return
		}

		healthMu.Lock()

		// stop health check go routine if one is running
		if healthCancel != nil {
			healthCancel()
		}

		ctx, cancel := context.WithCancel(context.Background())
		healthCancel = cancel

		healthMu.Unlock()

		go func(ctx context.Context, s Server) {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				ok, err := SyncService.CheckServerHealth(server)
				if err != nil || !ok {
					App.Event.Emit("serverHealth", "Offline")
				} else {
					App.Event.Emit("serverHealth", "Online")
				}

				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}
			}
		}(ctx, server)
	})

	// Run the application. This blocks until the application has been exited.
	err := App.Run()
	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
