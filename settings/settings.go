package settings

import (
	"io/ioutil"
	"log"
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
		SDLWindowWidth  int `yaml:"SDL_Window_Width"`
		SDLWindowHeight int `yaml:"SDL_Window_Height"`
	} `yaml:"AppWindow"`
	AppGui struct {
		GUIClearColor string `yaml:"guiClearColor"`
	} `yaml:"AppGui"`
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

	appConfig, err := ioutil.ReadFile("../resources/Kuplung_Settings.yaml")
	if err != nil {
		log.Fatalf("Settings error: %v", err)
	}

	err = yaml.Unmarshal(appConfig, &appSettings)
	if err != nil {
		log.Fatalf("Settings error: %v", err)
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

	err = ioutil.WriteFile("../resources/Kuplung_Settings.yaml", data, 0644)
	if err != nil {
		log.Fatalf("Settings save error: %v", err)
	}
}
