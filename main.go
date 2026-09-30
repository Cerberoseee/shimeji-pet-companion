package main

import (
	"embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"

	shimeji "shimeji-pet-companion/internal/shimeji"
	config "shimeji-pet-companion/internal/config"
	physics "shimeji-pet-companion/internal/physics"
)

var assets embed.FS

func init() {
	application.RegisterEvent[string]("time")
}

func main() {
	cfg := config.LoadConfig("config.yml")

	shimejiService := shimeji.NewShimejiService(cfg, nil)

	app := application.New(application.Options{
		Name:        "shimeji-pet-companion",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(shimejiService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	width := 200
	height := 200

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Shimeji Assistant",
		Width:            width,
		Height:           height,
		Frameless:        true,
		AlwaysOnTop:      true,
		BackgroundType:   application.BackgroundTypeTransparent,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		URL:              "/",
		Windows: application.WindowsWindow{
			DisableFramelessWindowDecorations: true,
			NonClientRegionSupport:            true,
		},
	})
	primaryScreen := app.Screen.GetPrimary()
	screenW, screenH := 1920, 1080

	if primaryScreen != nil {
		screenW, screenH = primaryScreen.Size.Width, primaryScreen.Size.Height
		window.SetPosition(
			primaryScreen.Size.Width-width-20,
			primaryScreen.Size.Height-height-20,
		)
	}
	shimeji.SetupTray(app, window)

	phyConfig := &physics.Config{
		TargetFPS:     24,
		Gravity:       900.0,
		Friction:      0.02,
		AirResistance: 0.85,
		Bounce:        0.28,
		FloorOffset:   48,
	}

	// Initialize physics engine
	engine := physics.NewEngine(app, window, screenW, screenH, width, height, phyConfig)
	shimejiService.SetPhysics(engine)
	engine.Start()

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
