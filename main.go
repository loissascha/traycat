package main

import (
	_ "embed"
	"fmt"
	"time"

	"fyne.io/systray"
)

//go:embed icon.png
var appIcon []byte

var cpuI *systray.MenuItem

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(appIcon)
	systray.SetTitle("Systray CAT")
	systray.SetTooltip("CPU x%")
	addCPUItem()
	addQuitItem()

	go func() {
		for {
			updateCPUPercentDisplay()
			time.Sleep(5 * time.Second)
		}
	}()
}

func addCPUItem() {
	cpuI = systray.AddMenuItem("CPU: 0%", "This is your CPU thing")
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

func updateCPUPercentDisplay() {
	systray.SetTooltip("CPU 1%")
	cpuI.SetTitle("CPU: 1%")
}

func onExit() {
	now := time.Now()
	fmt.Println("Exit at", now.String())
}
