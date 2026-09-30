package main

import (
	"embed"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"

	config "shimeji-pet-companion/internal/config"
	physics "shimeji-pet-companion/internal/physics"
	shimeji "shimeji-pet-companion/internal/shimeji"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[string]("time")
}

func main() {
	var cfg *config.Config
	if _, err := os.Stat("config.local.yml"); err == nil {
		cfg = config.LoadConfig("config.local.yml")
	} else {
		cfg = config.LoadConfig("config.yml")
	}

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
		Mac: application.MacWindow{
			DisableShadow: true,
			Backdrop:      application.MacBackdropTransparent,
		},
	})

	shimeji.SetupTray(app, window)

	phyConfig := &physics.Config{
		TargetFPS:     24,
		Gravity:       900.0,
		Friction:      0.02,
		AirResistance: 0.85,
		Bounce:        0.28,
		FloorOffset:   48,
	}

	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(_ *application.ApplicationEvent) {
		window.SetBackgroundColour(application.NewRGBA(0, 0, 0, 0))
	})

	var startPhysics sync.Once
	window.OnWindowEvent(events.Common.WindowRuntimeReady, func(_ *application.WindowEvent) {
		startPhysics.Do(func() {
			screenW, screenH := 1920, 1080
			if primaryScreen := app.Screen.GetPrimary(); primaryScreen != nil {
				screenW = primaryScreen.Size.Width
				if runtime.GOOS == "darwin" {
					screenH = primaryScreen.Size.Height - 30
				} else {
					screenH = primaryScreen.Size.Height
				}
			} else {
				log.Printf("primary screen unavailable after window runtime initialization; using %dx%d fallback", screenW, screenH)
			}

			engine := physics.NewEngine(app, window, screenW, screenH, width, height, phyConfig)
			shimejiService.SetPhysics(engine)
			engine.Start()
		})
	})

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
