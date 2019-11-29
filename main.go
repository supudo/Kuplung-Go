package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/supudo/Kuplung-Go/app"
	"github.com/supudo/Kuplung-Go/crash"
	"github.com/supudo/Kuplung-Go/settings"
)

var version string

func main() {
	scale := flag.Float64("scale", 1.0, "factor for scaling the UI (0.5 .. 10.0). 1080p displays should use default. 4K most likely 2.0.")
	fontFile := flag.String("fontfile", "", "Path to font file (.TTF) to use instead of the default font. Useful for HiDPI displays.")
	fontSize := flag.Float64("fontsize", 0.0, "Size of the font to use. If not specified, a default height will be used.")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	var kapp app.KuplungApp
	kapp.FontFile = *fontFile
	kapp.FontSize = float32(*fontSize)
	kapp.GuiScale = float32(*scale)
	if len(version) > 0 {
		kapp.Version = version
	} else {
		kapp.Version = fmt.Sprintf("(manual build %v)", time.Now().Format("2019-11-20"))
	}
	deferrer := make(chan func(), 100)

	versionInfo := "supudo.net - Kuplung-Go - " + kapp.Version
	defer crash.Handler(versionInfo)

	profileFin, err := initProfiling(*cpuprofile)
	if err != nil {
		settings.LogError("[main] Failed to start CPU profiling: %v\n", err)
	}
	defer profileFin()

	app.Run(kapp.InitializeKuplungWindow, versionInfo, 30.0, deferrer)
}

func initProfiling(filename string) (func(), error) {
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			return func() {}, err
		}
		err = pprof.StartCPUProfile(f)
		return func() { pprof.StopCPUProfile() }, err
	}
	return func() {}, nil
}
