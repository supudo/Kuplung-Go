package components

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// ComponentFileSaver ...
type ComponentFileSaver struct {
	showNewFolderModel    bool
	panelWidthFileOptions float32

	fileName      string
	newFolderName string
	currentFolder string

	positionX float32
	positionY float32
	width     float32
	height    float32
}

// NewComponentFileSaver ...
func NewComponentFileSaver() *ComponentFileSaver {
	sett := settings.GetSettings()
	return &ComponentFileSaver{
		positionX:             50,
		positionY:             50,
		width:                 sett.AppWindow.FileBrowserWidth,
		height:                sett.AppWindow.FileBrowserHeight,
		panelWidthFileOptions: 200.0,
		currentFolder:         sett.App.CurrentFolder,
		showNewFolderModel:    false,
	}
}

// Render ...
func (comp *ComponentFileSaver) Render(operation types.FileSaverOperation, open *bool) {
	sett := settings.GetSettings()

	imgui.SetNextWindowSizeV(imgui.Vec2{X: sett.AppWindow.FileBrowserWidth, Y: sett.AppWindow.FileBrowserHeight}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: comp.positionX, Y: comp.positionY}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	windowTitle := ""
	btnLabel := ""
	switch operation {
	case types.FileSaverOperationSaveScene:
		windowTitle = "Save Scene"
		btnLabel = "Save"
	case types.FileSaverOperationOpenScene:
		windowTitle = "Open Scene"
		btnLabel = "Open"
	case types.FileSaverOperationRenderer:
		windowTitle = "Render Scene"
	}

	if imgui.BeginV(windowTitle, open, 0) {
		imgui.Text(fmt.Sprintf("%s", filepath.Clean(comp.currentFolder)))
		imgui.Separator()

		// TODO: options
		// imgui.SameLine()

		imgui.BeginChild("scrolling")
		imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{X: 0, Y: 1})

		ww := float32(300)
		imgui.PushItemWidth(ww * 0.70)
		imgui.Text("File Name: ")
		imgui.InputText("", &comp.fileName)
		imgui.SameLineV(0, 10)
		if imgui.Button(btnLabel) {
			var file types.FBEntity

			file.IsFile = true
			file.Title = comp.fileName
			file.Path = comp.currentFolder + "/" + file.Title
			file.Extension = filepath.Ext(file.Title)
			file.ModifiedDate = ""
			file.Size = ""
			sett.App.CurrentFolder = comp.currentFolder
			settings.SaveSettings()

			switch operation {
			case types.FileSaverOperationSaveScene:
				_, _ = trigger.Fire(types.ActionFileSaverSaveScene, file)
			case types.FileSaverOperationOpenScene:
				_, _ = trigger.Fire(types.ActionFileSaverOpenScene, file)
			case types.FileSaverOperationRenderer:
				_, _ = trigger.Fire(types.ActionFileSaverRenderer, file)
			}
			*open = false
		}
		imgui.SameLineV(0, 10)
		if imgui.Button("New Folder") {
			comp.showNewFolderModel = true
		}
		imgui.PopItemWidth()
		imgui.Separator()

		if comp.showNewFolderModel {
			comp.modalNewFolder(ww)
		}

		imgui.ColumnsV(3, "fileColumns", 0)

		imgui.Separator()
		imgui.Text("File")
		imgui.NextColumn()
		imgui.Text("Size")
		imgui.NextColumn()
		imgui.Text("Last Modified")
		imgui.NextColumn()
		imgui.Separator()

		comp.drawFiles(open)

		imgui.ColumnsV(1, "", 0)

		imgui.Separator()
		imgui.Spacing()

		imgui.PopStyleVar()
		imgui.EndChild()

		imgui.End()
	}
}

func (comp *ComponentFileSaver) modalNewFolder(ww float32) {
	imgui.OpenPopup("New Folder")
	sett := settings.GetSettings()
	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(sett.AppWindow.SDLWindowWidth)/2 - 200, Y: float32(sett.AppWindow.SDLWindowHeight)/2 - 100}, imgui.ConditionAlways, imgui.Vec2{X: 0.5, Y: 0.5})
	imgui.SetNextWindowFocus()
	if imgui.BeginPopupModalV("New Folder", nil, imgui.WindowFlagsAlwaysAutoResize) {
		imgui.Text("Folder name:")

		if len(comp.newFolderName) == 0 {
			comp.newFolderName = "untitled"
		}
		imgui.PushItemWidth(ww)
		imgui.InputText("", &comp.newFolderName)
		imgui.PopItemWidth()

		if imgui.ButtonV("OK", imgui.Vec2{X: 200, Y: 0}) {
			newDir := comp.currentFolder + "/" + comp.newFolderName
			if _, err := os.Stat(newDir); os.IsNotExist(err) {
				if err := os.MkdirAll(newDir, 0755); err != nil {
					settings.LogError("[FileSaver] Cannot create new folder!")
				}
			}
			imgui.CloseCurrentPopup()
			comp.showNewFolderModel = false
			comp.newFolderName = ""
		}
		imgui.SameLineV(0, 20)
		if imgui.ButtonV("Cancel", imgui.Vec2{X: 140, Y: 0}) {
			imgui.CloseCurrentPopup()
			comp.showNewFolderModel = false
			comp.newFolderName = ""
		}

		imgui.EndPopup()
	}
}

func (comp *ComponentFileSaver) drawFiles(open *bool) {
	sett := settings.GetSettings()
	folderKeys, folderContents := comp.getFolderContents(sett.App.CurrentFolder)
	if runtime.GOOS == "windows" {
		// TODO: windows
		// if sett.CurrentDriveIndex != Settings::Instance()->Setting_SelectedDriveIndex) {
		// 	cFolder = Settings::Instance()->hddDriveList[Settings::Instance()->Setting_SelectedDriveIndex] + ":\\";
		// 	Settings::Instance()->Setting_CurrentDriveIndex = Settings::Instance()->Setting_SelectedDriveIndex;
		// 	this->currentFolder = cFolder;
		// }
	}

	i := int32(0)
	selected := int32(-1)
	for _, k := range folderKeys {
		entity := folderContents[k]
		if imgui.SelectableV(entity.Title, selected == i, imgui.SelectableFlagsSpanAllColumns, imgui.Vec2{X: 0, Y: 0}) {
			selected = i
			if entity.IsFile {
				comp.fileName = entity.Title
				sett.App.CurrentFolder = comp.currentFolder
				comp.currentFolder = sett.App.CurrentFolder
			} else {
				sett.App.CurrentFolder = entity.Path
				comp.currentFolder = sett.App.CurrentFolder
				comp.drawFiles(open)
			}
		}
		imgui.NextColumn()

		imgui.Text(entity.Size)
		imgui.NextColumn()
		imgui.Text(entity.ModifiedDate)
		imgui.NextColumn()

		i++
	}
}

func (comp *ComponentFileSaver) getFolderContents(filePath string) (folderKeys []string, folderContents map[string]*types.FBEntity) {
	currentPath := filepath.Clean(filePath)
	folderKeys = []string{}
	folderContents = make(map[string]*types.FBEntity)

	if settings.IsFolder(currentPath) {
		entity := &types.FBEntity{}
		entity.IsFile = false
		entity.Title = ".."
		entity.Path = filepath.Dir(currentPath)
		entity.Size = ""
		folderContents[entity.Path] = entity
		folderKeys = append(folderKeys, entity.Path)

		files, err := ioutil.ReadDir(currentPath)
		if err == nil {
			for _, f := range files {
				fext := filepath.Ext(f.Name())
				entity := &types.FBEntity{}
				entity.IsFile = !f.IsDir()
				if entity.IsFile {
					entity.Title = f.Name()
				} else {
					entity.Title = "<" + f.Name() + ">"
				}
				entity.Extension = fext
				entity.Path = currentPath + "/" + f.Name()
				if f.IsDir() {
					entity.Size = ""
				} else {
					entity.Size = settings.ConvertSize(f.Size())
				}
				entity.ModifiedDate = f.ModTime().Format("02-Jan-2006")
				if _, ok := folderContents[entity.Path]; !ok {
					folderContents[entity.Path] = entity
					folderKeys = append(folderKeys, entity.Path)
				}
			}
		}
	}

	sort.Strings(folderKeys)
	return folderKeys, folderContents
}
