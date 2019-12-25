package settings

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"

	"gopkg.in/yaml.v2"
)

// ApplicationSettings holds applications settings
type ApplicationSettings struct {
	App struct {
		ApplicationVersion string `yaml:"appVersion"`
		CurrentPath        string
		RendererType       uint32 `yaml:"RendererType"`
	} `yaml:"App"`
	AppWindow struct {
		SDLWindowWidth  float32 `yaml:"SDL_Window_Width"`
		SDLWindowHeight float32 `yaml:"SDL_Window_Height"`
		LogWidth        float32 `yaml:"Log_Width"`
		LogHeight       float32 `yaml:"Log_Height"`
	} `yaml:"AppWindow"`
	Rendering struct {
		FramesPerSecond float64 `yaml:"FramesPerSecond"`
	} `yaml:"Rendering"`
	AppGui struct {
		GUIClearColor []float32 `yaml:"guiClearColor"`
	} `yaml:"AppGui"`
	Consumption struct {
		ConsumptionIntervalCPU    int64 `yaml:"Consumption_Interval_CPU"`
		ConsumptionTimerCPU       int64
		ConsumptionCounterCPU     int64
		ConsumptionCPU            string
		ConsumptionIntervalMemory int64 `yaml:"Consumption_Interval_Memory"`
		ConsumptionTimerMemory    int64
		ConsumptionCounterMemory  int64
		ConsumptionMemory         string
	} `yaml:"Consumption"`
	MemSettings struct {
		QuitApplication bool

		LogBuffer      string
		LogBufferLimit int

		NbLastTime uint32
		NbFrames   uint32
		NbResult   float32

		TotalVertices  int32
		TotalIndices   int32
		TotalTriangles int32
		TotalFaces     int32
		TotalObjects   int32
	}
}

var instantiatedSettings *ApplicationSettings
var onceSettings sync.Once

// GetSettings singleton for our application settings
func GetSettings() *ApplicationSettings {
	onceSettings.Do(func() {
		as := InitSettings()
		instantiatedSettings = &as
	})
	return instantiatedSettings
}

// InitSettings will initialize application settings
func InitSettings() ApplicationSettings {
	var appSettings ApplicationSettings

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Settings error: %v", err)
	}

	if runtime.GOOS == "darwin" {
		appSettings.App.CurrentPath = dir + "/../Resources/resources/"
	} else if runtime.GOOS == "windows" {
		appSettings.App.CurrentPath = dir + "./"
	} else {
		// TODO: other platforms
		appSettings.App.CurrentPath = dir
	}

	appConfig, err := ioutil.ReadFile(appSettings.App.CurrentPath + "Kuplung_Settings.yaml")
	if err != nil {
		log.Fatalf("Settings error: %v", err)
	}

	err = yaml.Unmarshal(appConfig, &appSettings)
	if err != nil {
		log.Fatalf("Settings error: %v", err)
	}

	appSettings.MemSettings.QuitApplication = false
	appSettings.MemSettings.LogBuffer = ""
	appSettings.MemSettings.LogBufferLimit = 15360
	appSettings.MemSettings.NbLastTime = 0
	appSettings.MemSettings.NbFrames = 0
	appSettings.MemSettings.NbResult = 0.0
	appSettings.MemSettings.TotalVertices = 0
	appSettings.MemSettings.TotalIndices = 0
	appSettings.MemSettings.TotalTriangles = 0
	appSettings.MemSettings.TotalFaces = 0
	appSettings.MemSettings.TotalObjects = 0

	if appSettings.Rendering.FramesPerSecond == 0.0 {
		appSettings.Rendering.FramesPerSecond = 30.0
	}

	appSettings.Consumption.ConsumptionCPU = ""
	appSettings.Consumption.ConsumptionCounterCPU = 0
	appSettings.Consumption.ConsumptionTimerCPU = 0
	appSettings.Consumption.ConsumptionMemory = ""
	appSettings.Consumption.ConsumptionTimerMemory = 0
	appSettings.Consumption.ConsumptionCounterMemory = 0

	for idx, num := range appSettings.AppGui.GUIClearColor {
		appSettings.AppGui.GUIClearColor[idx] = num / 255.0
	}

	return appSettings
}

// SaveSettings will save the settings back to yaml file
func SaveSettings() {
	var sett = GetSettings()

	data, err := yaml.Marshal(&sett)
	if err != nil {
		log.Fatalf("Settings save error: %v", err)
	}

	err = ioutil.WriteFile(sett.App.CurrentPath+"Kuplung_Settings.yaml", data, 0644)
	if err != nil {
		log.Fatalf("Settings save error: %v", err)
	}
}
