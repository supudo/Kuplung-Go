package settings

import (
	"io/ioutil"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// ApplicationSettings holds applications settings
type ApplicationSettings struct {
	App struct {
		ApplicationVersion string `yaml:"appVersion"`
		CurrentPath        string `yaml:"currentFolder"`
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
	MemSettings struct {
		QuitApplication bool
		ZoomFactor      float32
		LogBuffer       string
		LogBufferLimit  int
	}
}

var instantiated *ApplicationSettings
var once sync.Once

// GetSettings singleton for our application settings
func GetSettings() *ApplicationSettings {
	once.Do(func() {
		as := InitSettings()
		instantiated = &as
	})
	return instantiated
}

// InitSettings will initialize application settings
func InitSettings() ApplicationSettings {
	var appSettings ApplicationSettings

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Settings error: %v", err)
	}

	appConfig, err := ioutil.ReadFile(dir + "/../Resources/resources/Kuplung_Settings.yaml")
	if err != nil {
		log.Fatalf("Settings error: %v", err)
	}

	err = yaml.Unmarshal(appConfig, &appSettings)
	if err != nil {
		log.Fatalf("Settings error: %v", err)
	}

	appSettings.App.CurrentPath = dir
	appSettings.MemSettings.QuitApplication = false
	appSettings.MemSettings.ZoomFactor = 45.0
	appSettings.MemSettings.LogBuffer = ""
	appSettings.MemSettings.LogBufferLimit = 15360

	if appSettings.Rendering.FramesPerSecond == 0.0 {
		appSettings.Rendering.FramesPerSecond = 30.0
	}

	for idx, num := range appSettings.AppGui.GUIClearColor {
		appSettings.AppGui.GUIClearColor[idx] = num / 255.0
	}

	return appSettings
}

// SaveSettings will save the settings back to yaml file
func SaveSettings() {
	var sett = GetSettings()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Settings error: %v", err)
	}

	data, err := yaml.Marshal(&sett)
	if err != nil {
		log.Fatalf("Settings save error: %v", err)
	}

	err = ioutil.WriteFile(dir+"/../Resources/resources/Kuplung_Settings.yaml", data, 0644)
	if err != nil {
		log.Fatalf("Settings save error: %v", err)
	}
}
