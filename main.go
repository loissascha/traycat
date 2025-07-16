package main

import (
	_ "embed"
	"flag"
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

//go:embed cats/light_cat_0.png
var lightCat0 []byte

//go:embed cats/light_cat_1.png
var lightCat1 []byte

//go:embed cats/light_cat_2.png
var lightCat2 []byte

//go:embed cats/light_cat_3.png
var lightCat3 []byte

//go:embed cats/light_cat_4.png
var lightCat4 []byte

var cpuI *systray.MenuItem

var theme *string
var lastAnimationId = 0
var ms = 150 * time.Millisecond

func main() {
	theme = flag.String("theme", "dark", "Use dark or light for either dark or light themed cat")
	flag.Parse()

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(darkCat0)
	systray.SetTitle("Systray CAT")
	addQuitItem()

	go func() {
		for {
			percent, err := cpu.Percent(0, false)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			totalPercent := percent[0]
			ms = 150 * time.Millisecond
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
			// fmt.Printf("CPU usage: %.2f%% MS: %v\n", totalPercent, ms)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			animateIcon()
			time.Sleep(ms)
		}
	}()
}

func animateIcon() {
	switch lastAnimationId {
	case 0:
		if *theme == "dark" {
			systray.SetIcon(darkCat1)
		} else {
			systray.SetIcon(lightCat1)
		}
		lastAnimationId = 1
		// fmt.Println("Update 1")
	case 1:
		if *theme == "dark" {
			systray.SetIcon(darkCat2)
		} else {
			systray.SetIcon(lightCat2)
		}
		lastAnimationId = 2
		// fmt.Println("Update 2")
	case 2:
		if *theme == "dark" {
			systray.SetIcon(darkCat3)
		} else {
			systray.SetIcon(lightCat3)
		}
		lastAnimationId = 3
		// fmt.Println("Update 3")
	case 3:
		if *theme == "dark" {
			systray.SetIcon(darkCat4)
		} else {
			systray.SetIcon(lightCat4)
		}
		lastAnimationId = 4
		// fmt.Println("Update 4")
	case 4:
		if *theme == "dark" {
			systray.SetIcon(darkCat0)
		} else {
			systray.SetIcon(lightCat0)
		}
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
