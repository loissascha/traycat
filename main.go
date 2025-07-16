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
var darkCatsPng embed.FS

//go:embed cats/dark_ico/*.ico
var darkCatsIco embed.FS

//go:embed cats/light_png/*.png
var lightCatsPng embed.FS

//go:embed cats/light_ico/*.ico
var lightCatsIco embed.FS

type Sprite struct {
	png []byte
	ico []byte
}

var catSprites map[int]Sprite

var cpuI *systray.MenuItem

var theme *string
var lastAnimationId = 0
var ms = 150 * time.Millisecond

func main() {
	theme = flag.String("theme", "dark", "Use dark or light themed cat")
	flag.Parse()

	catSprites = make(map[int]Sprite)

	if *theme == "dark" {
		for i := range 5 {
			cat, err := darkCatsPng.ReadFile(fmt.Sprintf("cats/dark_png/cat_%d.png", i))
			if err != nil {
				panic(fmt.Sprintf("no cat %d", i))
			}
			catIco, err := darkCatsIco.ReadFile(fmt.Sprintf("cats/dark_ico/cat_%d.ico", i))
			if err != nil {
				panic(fmt.Sprintf("no cat %d", i))
			}
			catSprites[i] = Sprite{
				png: cat,
				ico: catIco,
			}
		}
	} else {
		for i := range 5 {
			cat, err := lightCatsPng.ReadFile(fmt.Sprintf("cats/light_png/cat_%d.png", i))
			if err != nil {
				panic(fmt.Sprintf("no cat %d", i))
			}
			catIco, err := lightCatsIco.ReadFile(fmt.Sprintf("cats/light_ico/cat_%d.ico", i))
			if err != nil {
				panic(fmt.Sprintf("no cat %d", i))
			}
			catSprites[i] = Sprite{
				png: cat,
				ico: catIco,
			}
		}
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(catSprites[0].ico, catSprites[0].png)
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
	systray.SetTemplateIcon(p.ico, p.png)
}

func addQuitItem() {
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for range mQuit.ClickedCh {
			systray.Quit()
		}
	}()
}

func onExit() {
	fmt.Println("Off to chase a red dot...")
}
