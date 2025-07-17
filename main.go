package main

import (
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
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
var ciItem *systray.MenuItem

func main() {
	theme = flag.String("theme", "dark", "Use dark or light themed cat")
	install := flag.Bool("install", false, "Would you like to install the binary into your ~/.local/bin folder?")
	autostart := flag.Bool("autostart", false, "Would you like to add an autostart script for traycat?")
	flag.Parse()

	stopRunningInstances()

	if *install {
		installScript()
	}

	if *autostart {
		autostartScript()
	}

	if *install || *autostart {
		fmt.Println("Done.")
		return
	}

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

func stopRunningInstances() {
	fmt.Println("Stopping all running instances.")

	myPid := os.Getpid()

	processes, err := process.Processes()
	if err != nil {
		panic(err)
	}

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}
		if name == "traycat" && p.Pid != int32(myPid) {
			err := p.SendSignal(syscall.SIGKILL)
			if err != nil {
				fmt.Printf("Failed to kill process %d (%s): %v\n", p.Pid, name, err)
			} else {
				fmt.Printf("Process %d (%s) killed successfully.\n", p.Pid, name)
			}
		}
	}
}

func installScript() {
	fmt.Println("Installing traycat into your ~/.local/bin")

	_, err := os.Stat("./traycat")
	if os.IsNotExist(err) {
		panic("You are trying to install traycat but there is no file called `traycat` in you current directory.")
	} else if err != nil {
		panic(err)
	}

	source, err := os.Open("./traycat")
	if err != nil {
		panic(err)
	}
	defer source.Close()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	destination, err := os.Create(fmt.Sprintf("%s/.local/bin/traycat", homeDir))
	if err != nil {
		panic(err)
	}
	defer destination.Close()

	written, err := io.Copy(destination, source)
	if err != nil {
		panic(err)
	}
	fmt.Println("Written:", written)
}

func autostartScript() {
	_, err := os.Stat("~/.local/bin/traycat")
	if os.IsNotExist(err) {
		installScript()
	} else if err != nil {
		panic(err)
	}
}

func onReady() {
	systray.SetIcon(catSprites[0].png)
	ciItem = systray.AddMenuItem("CPU x%", "Current CPU usage")
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
			ciItem.SetTitle(fmt.Sprintf("CPU %.2f%%", totalPercent))
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
	systray.SetIcon(p.png)
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
