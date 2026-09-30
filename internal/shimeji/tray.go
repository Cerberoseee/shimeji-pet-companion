package shimeji

import (
	_ "embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Embed the raw icon bytes into the binary
//
//go:embed appicon.png
var trayIcon []byte

func SetupTray(app *application.App, window *application.WebviewWindow) *application.SystemTray {
	tray := app.SystemTray.New()
	tray.SetIcon(trayIcon)

	trayMenu := app.Menu.New()

	trayMenu.Add("Toggle DevTools").OnClick(func(ctx *application.Context) {
		window.OpenDevTools()
	})
	trayMenu.AddSeparator()
	trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	tray.SetMenu(trayMenu)

	return tray
}
