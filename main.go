package main

import (
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"time"

	"fyne.io/systray"
	"github.com/shirou/gopsutil/v3/cpu"
)

//go:embed cats/dark_png/*.png
var darkCats embed.FS

var catSprites map[int][]byte

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

	catSprites = make(map[int][]byte)

	if *theme == "dark" {
		for i := range 5 {
			cat, err := darkCats.ReadFile(fmt.Sprintf("cats/dark_png/cat_%d.png", i))
			if err != nil {
				panic(fmt.Sprintf("no cat %d", i))
			}
			catSprites[i] = cat
		}
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(catSprites[0])
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
	lastAnimationId++
	p, ok := catSprites[lastAnimationId]
	if !ok {
		p = catSprites[0]
		lastAnimationId = 0
	}
	systray.SetIcon(p)
}

func addQuitItem() {
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
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
