package components

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// ComponentImport ...
type ComponentImport struct {
	positionX, positionY int32
	width, height        int32
	panelWidthOptions    float32
	panelWidthOptionsMin float32

	SettingForward, SettingUp int32

	currentFolder string

	formats  []string
	forwards []string
	ups      []string
	parsers  []string

	dialogImportType types.ImportExportFormat
}

// NewComponentImport ...
func NewComponentImport() *ComponentImport {
	sett := settings.GetSettings()
	comp := &ComponentImport{}
	comp.panelWidthOptions = 200.0
	comp.panelWidthOptionsMin = 200.0
	comp.SettingForward = 2
	comp.SettingUp = 4
	comp.currentFolder = sett.App.CurrentFolder
	comp.formats = []string{
		"Wavefront OBJ",
		"glTF",
		"STereoLithography STL",
		"Stanford PLY"}
	comp.forwards = []string{
		"-X Forward",
		"-Y Forward",
		"-Z Forward",
		"X Forward",
		"Y Forward",
		"Z Forward"}
	comp.ups = []string{
		"-X Up",
		"-Y Up",
		"-Z Up",
		"X Up",
		"Y Up",
		"Z Up"}
	comp.parsers = []string{
		"Kuplung"}
	return comp
}

// Render ...
func (comp *ComponentImport) Render(open *bool, dialogImportType *types.ImportExportFormat) {
	comp.dialogImportType = *dialogImportType
	sett := settings.GetSettings()

	imgui.SetNextWindowSizeV(imgui.Vec2{X: sett.AppWindow.FileBrowserWidth, Y: sett.AppWindow.FileBrowserHeight}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: 40, Y: 40}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	windowTitle := "Import "
	switch *dialogImportType {
	case types.ImportExportFormatOBJ:
		windowTitle += "Wavefront OBJ file"
	case types.ImportExportFormatSTL:
		windowTitle += "STereoLithography STL file"
	case types.ImportExportFormatPLY:
		windowTitle += "Stanford PLY file"
	case types.ImportExportFormatGLTF:
		windowTitle += "glTF file"
	}

	if imgui.BeginV(windowTitle, open, 0) {
		imgui.Text(fmt.Sprintf("%s", filepath.Clean(comp.currentFolder)))
		imgui.Separator()

		imgui.BeginChildV("OptionsPanel", imgui.Vec2{X: comp.panelWidthOptions, Y: 0}, true, 0)

		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
		imgui.Text("Options")
		imgui.PopStyleColorV(1)
		imgui.Separator()
		imgui.PushItemWidth(-1.0)
		imgui.Text("Kuplung File Format")
		if imgui.BeginCombo("##982", comp.formats[int32(*dialogImportType)]) {
			var i int32
			for i = 0; i < int32(len(comp.formats)); i++ {
				if imgui.SelectableV(comp.formats[i], (types.ImportExportFormat(i) == *dialogImportType), 0, imgui.Vec2{X: 0, Y: 0}) {
					*dialogImportType = types.ImportExportFormat(i)
				}
			}
			imgui.EndCombo()
		}
		imgui.PopItemWidth()
		imgui.Separator()

		if runtime.GOOS == "windows" {
			imgui.Separator()
			imgui.Text("Select Drive")
			// TODO: show drives
			imgui.Separator()
		}

		imgui.Separator()
		imgui.PushItemWidth(-1.0)
		imgui.Text("Forward")
		var i int32
		if imgui.BeginCombo("##987", comp.forwards[comp.SettingForward]) {
			for i = 0; i < int32(len(comp.forwards)); i++ {
				if imgui.SelectableV(comp.forwards[i], (i == comp.SettingForward), 0, imgui.Vec2{X: 0, Y: 0}) {
					comp.SettingForward = i
				}
			}
			imgui.EndCombo()
		}
		imgui.Separator()
		imgui.Text("Up")
		if imgui.BeginCombo("##988", comp.ups[comp.SettingUp]) {
			for i = 0; i < int32(len(comp.ups)); i++ {
				if imgui.SelectableV(comp.ups[i], (i == comp.SettingUp), 0, imgui.Vec2{X: 0, Y: 0}) {
					comp.SettingUp = i
				}
			}
			imgui.EndCombo()
		}
		imgui.Separator()
		if imgui.ButtonV("From Blender", imgui.Vec2{X: -1.0, Y: 0.0}) {
			comp.SettingForward = 2
			comp.SettingUp = 4
		}
		imgui.Separator()
		imgui.Text("Parser:")
		// TODO: cuda parsers
		if imgui.BeginCombo("##989", comp.parsers[sett.MemSettings.ModelFileParser]) {
			for i = 0; i < int32(len(comp.parsers)); i++ {
				if imgui.SelectableV(comp.parsers[i], (i == sett.MemSettings.ModelFileParser), 0, imgui.Vec2{X: 0, Y: 0}) {
					sett.MemSettings.ModelFileParser = i
				}
			}
			imgui.EndCombo()
		}
		imgui.PopItemWidth()
		imgui.EndChild()

		imgui.SameLine()

		// sc := float32(1.0 / 255.0)
		// imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 89.0 * sc, Y: 91.0 * sc, Z: 94 * sc, W: 1})
		// imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: 119.0 * sc, Y: 122.0 * sc, Z: 124.0 * sc, W: 1})
		// imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: .0, Y: .0, Z: .0, W: 1})
		// imgui.ButtonV("###splitterOptionsImport", imgui.Vec2{X: -1, Y: 4})
		// imgui.PopStyleColorV(3)
		// // TODO: get mouse delta up/down
		// // 	this->panelWidth_Options += imgui.GetIO().MouseDelta.x;
		// // 	if (this->panelWidth_Options < this->panelWidth_OptionsMin)
		// // 		this->panelWidth_Options = this->panelWidth_OptionsMin;
		// // }
		// if imgui.IsItemHovered() {
		// 	imgui.SetMouseCursor(imgui.MouseCursorResizeNS)
		// } else {
		// 	imgui.SetMouseCursor(imgui.MouseCursorNone)
		// }

		// imgui.SameLine()

		// folder browser
		imgui.BeginChild("scrolling")
		imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{X: 0, Y: 1})

		imgui.ColumnsV(3, "fileColumns", 0)

		imgui.Separator()
		imgui.Text("File")
		imgui.NextColumn()
		imgui.Text("Size")
		imgui.NextColumn()
		imgui.Text("Last Modified")
		imgui.NextColumn()
		imgui.Separator()

		comp.drawFiles(dialogImportType, open)

		imgui.ColumnsV(1, "", 0)

		imgui.Separator()
		imgui.Spacing()

		imgui.PopStyleVar()
		imgui.EndChild()

		imgui.End()
	}
}

func (comp *ComponentImport) drawFiles(dialogImportType *types.ImportExportFormat, open *bool) {
	sett := settings.GetSettings()
	folderKeys, folderContents := comp.getFolderContents(dialogImportType, sett.App.CurrentFolder)
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
				var setts []string
				setts = append(setts, fmt.Sprintf("%v", comp.SettingForward))
				setts = append(setts, fmt.Sprintf("%v", comp.SettingUp))
				_, _ = trigger.Fire(types.ActionFileImport, entity, setts, *dialogImportType)

				sett.App.CurrentFolder = comp.currentFolder
				comp.currentFolder = sett.App.CurrentFolder
				settings.SaveSettings()
				*open = false
			} else {
				sett.App.CurrentFolder = entity.Path
				comp.currentFolder = sett.App.CurrentFolder
				comp.drawFiles(dialogImportType, open)
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

func (comp *ComponentImport) getFolderContents(dialogImportType *types.ImportExportFormat, filePath string) (folderKeys []string, folderContents map[string]*types.FBEntity) {
	currentPath := filepath.Clean(filePath)
	folderKeys = []string{}
	folderContents = make(map[string]*types.FBEntity)

	if comp.isFolder(currentPath) {
		entity := &types.FBEntity{}
		entity.IsFile = false
		entity.Title = ".."
		entity.Path = filepath.Dir(currentPath)
		entity.Size = ""
		folderContents[entity.Path] = entity
		folderKeys = append(folderKeys, entity.Path)

		files, err := ioutil.ReadDir(currentPath)
		if err == nil {
			isAllowedFileExtension := false
			for _, f := range files {
				fext := filepath.Ext(f.Name())
				if *dialogImportType != types.ImportExportFormatUNDEFINED {
					switch *dialogImportType {
					case types.ImportExportFormatOBJ:
						isAllowedFileExtension = fext == ".obj"
					case types.ImportExportFormatGLTF:
						isAllowedFileExtension = fext == ".gltf"
					case types.ImportExportFormatPLY:
						isAllowedFileExtension = fext == ".ply"
					case types.ImportExportFormatSTL:
						isAllowedFileExtension = fext == ".stl"
					}
				} else {
					isAllowedFileExtension = true
				}
				if isAllowedFileExtension || f.IsDir() {
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
						entity.Size = comp.convertSize(f.Size())
					}
					entity.ModifiedDate = f.ModTime().Format("02-Jan-2006")
					if _, ok := folderContents[entity.Path]; !ok {
						folderContents[entity.Path] = entity
						folderKeys = append(folderKeys, entity.Path)
					}
				}
			}
		}
	}

	sort.Strings(folderKeys)
	return folderKeys, folderContents
}

func (comp *ComponentImport) convertSize(size int64) string {
	sizes := []string{"B", "KB", "MB", "GB"}
	div := int32(0)
	rem := float32(0)

	for size >= 1024 && div < int32(len(sizes)) {
		rem = float32(int32(size) % 1024)
		div++
		size /= 1024
	}

	sized := float32(size) + rem/1024.0
	result := fmt.Sprintf("%.2f %s", (comp.roundOff(sized)), sizes[div])
	return result
}

func (comp *ComponentImport) roundOff(n float32) float32 {
	d := n * 100.0
	i := d + 0.5
	d = i / 100.0
	return d
}

func (comp *ComponentImport) isHidden(p string) bool {
	name := filepath.Base(p)
	if name == ".." || name == "." || strings.HasPrefix(name, ".") {
		return true
	}
	return false
}

func (comp *ComponentImport) isFolder(path string) bool {
	fileInfo, _ := os.Stat(path)
	return fileInfo.IsDir()
}
