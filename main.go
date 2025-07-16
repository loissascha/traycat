package main

import (
	_ "embed"
	"fmt"
	"time"

	"fyne.io/systray"
	"github.com/shirou/gopsutil/v3/cpu"
)

//go:embed cats/dark_cat_0.png
var darkCat0 []byte

//go:embed cats/dark_cat_1.png
var darkCat1 []byte

//go:embed cats/dark_cat_2.png
var darkCat2 []byte

//go:embed cats/dark_cat_3.png
var darkCat3 []byte

//go:embed cats/dark_cat_4.png
var darkCat4 []byte

var cpuI *systray.MenuItem

var lastAnimationId = 0

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(darkCat0)
	systray.SetTitle("Systray CAT")
	// systray.SetTooltip("CPU x%")
	addQuitItem()

	go func() {
		for {
			animateIcon()
			percent, err := cpu.Percent(0, false)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			totalPercent := percent[0]
			ms := 150 * time.Millisecond
			if totalPercent > 10 {
				ms = 120 * time.Millisecond
			}
			if totalPercent > 20 {
				ms = 100 * time.Millisecond
			}
			if totalPercent > 30 {
				ms = 80 * time.Millisecond
			}
			if totalPercent > 40 {
				ms = 60 * time.Millisecond
			}
			if totalPercent > 50 {
				ms = 40 * time.Millisecond
			}

			// systray.SetTooltip(fmt.Sprintf("CPU %.2f%%", totalPercent))
			// fmt.Printf("Total CPU Usage: %.2f%%\n", percent[0])
			// fmt.Println("Sleep time:", ms)
			time.Sleep(ms)
		}
	}()

	// go func() {
	// 	for {
	// 		updateCPUPercentDisplay()
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()
}

func animateIcon() {
	switch lastAnimationId {
	case 0:
		systray.SetIcon(darkCat1)
		lastAnimationId = 1
		// fmt.Println("Update 1")
	case 1:
		systray.SetIcon(darkCat2)
		lastAnimationId = 2
		// fmt.Println("Update 2")
	case 2:
		systray.SetIcon(darkCat3)
		lastAnimationId = 3
		// fmt.Println("Update 3")
	case 3:
		systray.SetIcon(darkCat4)
		lastAnimationId = 4
		// fmt.Println("Update 4")
	case 4:
		systray.SetIcon(darkCat0)
		lastAnimationId = 0
		// fmt.Println("Update 0")
	}
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
