package settings

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/supudo/Kuplung-Go/types"
)

// ReadFile ...
func ReadFile(filename string, terminated bool) string {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		LogWarn("[Settings] Can't load shader source for %v", filename)
		return ""
	}
	if terminated {
		return string(source) + "\x00"
	}
	return string(source) + " "
}

// SaveRecentFilesImported ...
func SaveRecentFilesImported(recentFilesImported []*types.FBEntity) {
	sett := GetSettings()
	nlDelimiter := ""
	if runtime.GOOS == "darwin" {
		nlDelimiter = "\n"
	}
	recentFilesLines := "# Recent Imported Files list" + nlDelimiter + nlDelimiter
	for i := 0; i < len(recentFilesImported); i++ {
		fileEntity := recentFilesImported[i]
		recentFilesLines += "# File" + nlDelimiter
		if len(fileEntity.Title) == 0 {
			recentFilesLines += "-" + nlDelimiter
		} else {
			recentFilesLines += fileEntity.Title + nlDelimiter
		}
		if len(fileEntity.Path) == 0 {
			recentFilesLines += "-" + nlDelimiter
		} else {
			recentFilesLines += fileEntity.Path + nlDelimiter
		}
		recentFilesLines += nlDelimiter
	}
	SaveStringToFile(recentFilesLines, sett.App.AppFolder+"Kuplung_RecentFiles.ini", "")
}

// SaveStringToFile ...
func SaveStringToFile(fileContents, filepath, message string) {
	var f *os.File
	var err error
	if _, _ = os.Stat(filepath); os.IsNotExist(err) {
		f, err = os.Create(filepath)
		if err != nil {
			LogError("[Settings] [%v] Can't create file : %v!", message, filepath)
		}
		defer f.Close()
	} else {
		f, err = os.OpenFile(filepath, os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			LogError("[Settings] [%v] Can't open file : %v!", message, filepath)
		}
	}
	_, err = f.WriteString(fileContents)
	if err != nil {
		LogError("[Settings] [%v] Can't save file : %v!", message, filepath)
	}
	f.Sync()
}

// ConvertSize ...
func ConvertSize(size int64) string {
	sizes := []string{"B", "KB", "MB", "GB"}
	div := int32(0)
	rem := float32(0)

	for size >= 1024 && div < int32(len(sizes)) {
		rem = float32(int32(size) % 1024)
		div++
		size /= 1024
	}

	sized := float32(size) + rem/1024.0
	result := fmt.Sprintf("%.2f %s", (roundOff(sized)), sizes[div])
	return result
}

// IsFolder ...
func IsFolder(path string) bool {
	fileInfo, _ := os.Stat(path)
	return fileInfo.IsDir()
}

func roundOff(n float32) float32 {
	d := n * 100.0
	i := d + 0.5
	d = i / 100.0
	return d
}

func isHidden(p string) bool {
	name := filepath.Base(p)
	if name == ".." || name == "." || strings.HasPrefix(name, ".") {
		return true
	}
	return false
}
