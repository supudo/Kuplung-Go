package settings

import "sync"

// ApplicationSettings holds applications settings
type ApplicationSettings struct {
	ApplicationVersion string
}

var instantiated *ApplicationSettings
var once sync.Once

// GetSettings singleton for our application settings
func GetSettings() *ApplicationSettings {
	once.Do(func() {
		instantiated = &ApplicationSettings{}
		InitSettings()
	})
	return instantiated
}

// InitSettings will initialize application settings
func InitSettings() ApplicationSettings {
	var appSettings ApplicationSettings
	appSettings.ApplicationVersion = "1.0"
	return appSettings
}
