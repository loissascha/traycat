package main

import (
	_ "embed"
	"fmt"
	"time"

	"fyne.io/systray"
)

//go:embed icon.png
var appIcon []byte

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(appIcon)
	systray.SetTitle("Systray CAT")
	systray.SetTooltip("CPU x%")
	addQuitItem()

	// go func() {
	// 	for {
	// 		updateCPUPercent()
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()
}

func addCPUItem() {
	cpuI := systray.AddMenuItem("CPU: 1%", "This is your CPU thing")
	cpuI.SetIcon(appIcon)
}

func addQuitItem() {
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	// Sets the icon of a menu item.
	// mQuit.SetIcon(icon.Data)
	go func() {
		for range mQuit.ClickedCh {
			fmt.Println("Requesting quit")
			systray.Quit()
		}
	}()
}

func updateCPUPercent() {
	systray.ResetMenu()
	addCPUItem()
	addQuitItem()
}

func onExit() {
	now := time.Now()
	fmt.Println("Exit at", now.String())
}
