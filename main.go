package main

import (
	"fmt"
	"time"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Systray CAT")
	systray.SetTooltip("CPU x%")
	addQuitItem()
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

func onExit() {
	now := time.Now()
	fmt.Println("Exit at", now.String())
}
