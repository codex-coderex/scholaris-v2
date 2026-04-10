package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// embed

//go:embed all:frontend/dist
var assets embed.FS

// main

func main() {
	app := application.New(application.Options{
		Name:        "Scholaris V2",
		Description: "Student Information System",
		Services: []application.Service{
			application.NewService(NewAppService()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Scholaris V2",
		Width:            1280,
		Height:           800,
		MinWidth:         1280,
		MinHeight:        800,
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
