package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/supudo/Kuplung-Go/app"
	"github.com/supudo/Kuplung-Go/settings"
)

var version string

func main() {
	mbv := time.Now().Format("2019-11-20")
	var kapp app.KuplungApp
	if len(version) > 0 {
		kapp.Version = version
	} else {
		kapp.Version = fmt.Sprintf("(manual build %v)", mbv)
	}
	deferrer := make(chan func(), 100)

	versionInfo := "Kuplung-Go - " + kapp.Version

	profileFin, err := initProfiling(mbv)
	if err != nil {
		settings.LogError("[main] Failed to start CPU profiling: %v\n", err)
	}
	defer profileFin()

	app.Run(kapp.InitializeKuplungWindow, versionInfo, 30.0, deferrer)
}

func initProfiling(mbv string) (func(), error) {
	filename := "Kuplung_CPUProfile" + mbv + ".log"
	f, err := os.Create(filename)
	if err != nil {
		return func() {}, err
	}
	err = pprof.StartCPUProfile(f)
	return func() { pprof.StopCPUProfile() }, err
}
