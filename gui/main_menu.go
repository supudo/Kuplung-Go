package gui

import (
	"fmt"
	"os"
	"runtime"

	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/gui/fonts"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
	"github.com/supudo/Kuplung-Go/utilities"
)

// DrawMainMenu ...
func (context *Context) DrawMainMenu() {
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()

	// Main Menu
	imgui.BeginMainMenuBar()

	if imgui.BeginMenu("File") {
		// TODO: add FA icons
		if imgui.MenuItem(fmt.Sprintf("%c New", fonts.FA_ICON_FILE_O)) {
			_, _ = trigger.Fire(types.ActionGuiActionFileNew)
		}
		if imgui.MenuItem(fmt.Sprintf("%c Open ...", fonts.FA_ICON_FOLDER_OPEN_O)) {
			context.GuiVars.showOpenDialog = true
		}
		if imgui.BeginMenu(fmt.Sprintf("%c Open Recent", fonts.FA_ICON_FILES_O)) {
			// if (this->recentFiles.size() == 0)
			// 	imgui.MenuItem("No recent files", nil, false, false);
			// else {
			// 	for (size_t i = 0; i < this->recentFiles.size(); i++) {
			// 		FBEntity file = this->recentFiles[i];
			// 		if (imgui.MenuItem(file.title.c_str(), nil, false, true)) {
			// 			if (boost::filesystem::exists(file.path))
			// 				this->funcOpenScene(file);
			// 			else
			// 				this->showRecentFileDoesntExists = true;
			// 		}
			// 	}
			// 	imgui.Separator();
			// 	if (imgui.MenuItem("Clear recent files", nil, false))
			// 		this->recentFilesClear();
			// }
			imgui.EndMenu()
		}

		if imgui.MenuItem(fmt.Sprintf("%c Save ...", fonts.FA_ICON_FLOPPY_O)) {
			context.GuiVars.showSaveDialog = true
		}

		imgui.Separator()

		if imgui.BeginMenu("   Import") {
			if imgui.MenuItemV("Wavefront (.OBJ)", "", context.GuiVars.showImporterFile, true) {
				context.GuiVars.showImporterFile = true
				context.GuiVars.dialogImportType = types.ImportExportFormatOBJ
			}
			if imgui.MenuItemV("glTF (.gltf)", "", context.GuiVars.showImporterFile, true) {
				context.GuiVars.showImporterFile = true
				context.GuiVars.dialogImportType = types.ImportExportFormatGLTF
			}
			if imgui.MenuItemV("STereoLithography (.STL)", "", context.GuiVars.showImporterFile, true) {
				context.GuiVars.showImporterFile = true
				context.GuiVars.dialogImportType = types.ImportExportFormatSTL
			}
			if imgui.MenuItemV("Stanford (.PLY)", "", context.GuiVars.showImporterFile, true) {
				context.GuiVars.showImporterFile = true
				context.GuiVars.dialogImportType = types.ImportExportFormatPLY
			}
			if imgui.BeginMenu("Assimp...") {
				// for (size_t a = 0; a < Settings::Instance()->AssimpSupportedFormats_Import.size(); a++) {
				// 	SupportedAssimpFormat format = Settings::Instance()->AssimpSupportedFormats_Import[a];
				// 	std::string f = std::string(format.description) + " (" + std::string(format.fileExtension) + ")";
				// 	if (imgui.MenuItem(f.c_str(), nil, &this->showImporterFile))
				// 		this->dialogImportType_Assimp = static_cast<int>(a);
				// }
				imgui.EndMenu()
			}
			imgui.EndMenu()
		}

		if imgui.BeginMenu(fmt.Sprintf("%c Import Recent", fonts.FA_ICON_FILES_O)) {
			if len(context.GuiVars.recentFilesImported) == 0 {
				imgui.MenuItem("No recent files")
			} else {
				for i := 0; i < len(context.GuiVars.recentFilesImported); i++ {
					file := context.GuiVars.recentFilesImported[i]
					if imgui.MenuItem(file.Title) {
						if _, err := os.Stat(file.Path); !os.IsNotExist(err) {
							var setts []string
							setts = append(setts, "2")
							setts = append(setts, "4")
							_, _ = trigger.Fire(types.ActionFileImport, file, setts, context.GuiVars.dialogImportType)
						} else {
							context.GuiVars.showRecentFileImportedDoesntExists = true
						}
					}
				}
				imgui.Separator()
				if imgui.MenuItem("Clear recent files") {
					context.recentFilesClearImported()
				}
			}
			imgui.EndMenu()
		}

		if imgui.BeginMenu("   Export") {
			if imgui.MenuItemV("Wavefront (.OBJ)", "", context.GuiVars.showExporterFile, true) {
				context.GuiVars.showExporterFile = true
				context.GuiVars.dialogExportType = types.ImportExportFormatOBJ
			}
			if imgui.MenuItemV("glTF (.gltf)", "", context.GuiVars.showExporterFile, true) {
				context.GuiVars.showExporterFile = true
				context.GuiVars.dialogExportType = types.ImportExportFormatGLTF
			}
			if imgui.MenuItemV("STereoLithography (.stl)", "", context.GuiVars.showExporterFile, true) {
				context.GuiVars.showExporterFile = true
				context.GuiVars.dialogExportType = types.ImportExportFormatSTL
			}
			if imgui.MenuItemV("Stanford PLY (.ply)", "", context.GuiVars.showExporterFile, true) {
				context.GuiVars.showExporterFile = true
				context.GuiVars.dialogExportType = types.ImportExportFormatPLY
			}
			if imgui.BeginMenu("Assimp...") {
				// for (size_t a = 0; a < Settings::Instance()->AssimpSupportedFormats_Export.size(); a++) {
				// 	SupportedAssimpFormat format = Settings::Instance()->AssimpSupportedFormats_Export[a];
				// 	std::string f = std::string(format.description) + " (" + std::string(format.fileExtension) + ")";
				// 	if (imgui.MenuItem(f.c_str(), nil, &this->showExporterFile))
				// 		this->dialogExportType_Assimp = static_cast<int>(a);
				// }
				imgui.EndMenu()
			}
			imgui.EndMenu()
		}

		imgui.Separator()

		quitShortcut := ""
		if runtime.GOOS == "darwin" {
			quitShortcut = "Cmd+Q"
		} else if runtime.GOOS == "windows" {
			quitShortcut = "Alt+F4"
		}
		if imgui.MenuItemV(fmt.Sprintf("%c Quit", fonts.FA_ICON_POWER_OFF), quitShortcut, false, true) {
			os.Exit(3)
		}
		imgui.EndMenu()
	}

	if imgui.BeginMenu("Scene") {
		if imgui.BeginMenu(fmt.Sprintf("%c Add Light", fonts.FA_ICON_LIGHTBULB_O)) {
			if imgui.MenuItem("Directional (Sun)") {
				_, _ = trigger.Fire(types.ActionGuiAddLight, types.LightSourceTypeDirectional)
			}
			if imgui.MenuItem("Point (Light bulb)") {
				_, _ = trigger.Fire(types.ActionGuiAddLight, types.LightSourceTypePoint)
			}
			if imgui.MenuItem("Spot (Flashlight)") {
				_, _ = trigger.Fire(types.ActionGuiAddLight, types.LightSourceTypeSpot)
			}
			imgui.EndMenu()
		}
		imgui.Separator()
		if imgui.BeginMenu(fmt.Sprintf("%c Scene Rendering", fonts.FA_ICON_CERTIFICATE)) {
			if imgui.MenuItemV("Solid", "", rsett.General.SelectedViewModelSkin == types.ViewModelSkinSolid, true) {
				rsett.General.SelectedViewModelSkin = types.ViewModelSkinSolid
			}
			if imgui.MenuItemV("Material", "", rsett.General.SelectedViewModelSkin == types.ViewModelSkinMaterial, true) {
				rsett.General.SelectedViewModelSkin = types.ViewModelSkinMaterial
			}
			if imgui.MenuItemV("Texture", "", rsett.General.SelectedViewModelSkin == types.ViewModelSkinTexture, true) {
				rsett.General.SelectedViewModelSkin = types.ViewModelSkinTexture
			}
			if imgui.MenuItemV("Wireframe", "", rsett.General.SelectedViewModelSkin == types.ViewModelSkinWireframe, true) {
				rsett.General.SelectedViewModelSkin = types.ViewModelSkinWireframe
			}
			if imgui.MenuItemV("Rendered", "", rsett.General.SelectedViewModelSkin == types.ViewModelSkinRendered, true) {
				rsett.General.SelectedViewModelSkin = types.ViewModelSkinRendered
			}
			imgui.Separator()
			imgui.MenuItemV("Render - Depth", "", rsett.General.RenderingDepth, true)
			imgui.EndMenu()
		}
		imgui.Separator()
		imgui.MenuItemV(fmt.Sprintf("%c Render Image", fonts.FA_ICON_FILE_IMAGE_O), "", context.GuiVars.showImageSave, true)
		imgui.MenuItemV(fmt.Sprintf("%c Renderer UI", fonts.FA_ICON_CUBES), "", context.GuiVars.showRendererUI, true)

		imgui.EndMenu()
	}

	if imgui.BeginMenu("View") {
		lbl := fmt.Sprintf("%c GUI Controls", fonts.FA_ICON_TOGGLE_ON)
		if context.GuiVars.showModels {
			lbl = fmt.Sprintf("%c GUI Controls", fonts.FA_ICON_TOGGLE_OFF)
		}
		if imgui.MenuItem(lbl) {
			context.GuiVars.showModels = !context.GuiVars.showModels
		}
		lbl = fmt.Sprintf("%c Controls", fonts.FA_ICON_TOGGLE_ON)
		if context.GuiVars.showControls {
			lbl = fmt.Sprintf("%c Controls", fonts.FA_ICON_TOGGLE_OFF)
		}
		if imgui.MenuItem(lbl) {
			context.GuiVars.showControls = !context.GuiVars.showControls
		}
		lbl = fmt.Sprintf("%c Show Visual Artefacts", fonts.FA_ICON_TOGGLE_ON)
		if rsett.General.ShowAllVisualArtefacts {
			lbl = fmt.Sprintf("%c Hide Visual Artefacts", fonts.FA_ICON_TOGGLE_OFF)
		}
		if imgui.MenuItem(lbl) {
			rsett.General.ShowAllVisualArtefacts = !rsett.General.ShowAllVisualArtefacts
		}
		imgui.Separator()
		if imgui.MenuItem(fmt.Sprintf("%c Show Log Window", fonts.FA_ICON_BUG)) {
			context.GuiVars.showLog = !context.GuiVars.showLog
		}
		if sett.App.RendererType == types.InAppRendererTypeForward {
			if imgui.MenuItem(fmt.Sprintf("%c IDE", fonts.FA_ICON_PENCIL)) {
				context.GuiVars.showKuplungIDE = !context.GuiVars.showKuplungIDE
			}
		}
		if imgui.MenuItem(fmt.Sprintf("%c Screenshot", fonts.FA_ICON_DESKTOP)) {
			context.GuiVars.showScreenshotWindow = !context.GuiVars.showScreenshotWindow
		}
		if imgui.MenuItem(fmt.Sprintf("%c Scene Statistics", fonts.FA_ICON_TACHOMETER)) {
			context.GuiVars.showSceneStats = !context.GuiVars.showSceneStats
		}
		if imgui.MenuItem(fmt.Sprintf("%c Structured Volumetric Sampling", fonts.FA_ICON_PAPER_PLANE_O)) {
			context.GuiVars.showSVS = !context.GuiVars.showSVS
		}
		if imgui.MenuItem(fmt.Sprintf("%c Shadertoy", fonts.FA_ICON_PAPER_PLANE_O)) {
			context.GuiVars.showShadertoy = !context.GuiVars.showShadertoy
		}
		imgui.Separator()
		if imgui.MenuItem(fmt.Sprintf("%c Options", fonts.FA_ICON_COG)) {
			context.GuiVars.showOptions = !context.GuiVars.showOptions
		}
		imgui.Separator()
		if imgui.MenuItem(fmt.Sprintf("%c Cube", fonts.FA_ICON_CUBE)) {
			rsett.General.ShowCube = !rsett.General.ShowCube
		}
		imgui.EndMenu()
	}

	if imgui.BeginMenu("Help") {
		if imgui.MenuItem(fmt.Sprintf("%c Metrics", fonts.FA_ICON_INFO)) {
			context.GuiVars.showMetrics = !context.GuiVars.showMetrics
		}
		if imgui.MenuItem(fmt.Sprintf("%c About ImGui", fonts.FA_ICON_INFO_CIRCLE)) {
			context.GuiVars.showAboutImGui = !context.GuiVars.showAboutImGui
		}
		if imgui.MenuItem(fmt.Sprintf("%c About Kuplung", fonts.FA_ICON_INFO_CIRCLE)) {
			context.GuiVars.showAboutKuplung = !context.GuiVars.showAboutKuplung
		}
		imgui.Separator()
		if imgui.MenuItem("   ImGui Demo Window") {
			context.GuiVars.showDemoWindow = !context.GuiVars.showDemoWindow
		}
		imgui.EndMenu()
	}

	imgui.Text(fmt.Sprintf(" | [%.4f ms/frame] %d objs, %d verts, %d indices (%d tris, %d faces) | %v", sett.MemSettings.NbResult, sett.MemSettings.TotalObjects, sett.MemSettings.TotalVertices, sett.MemSettings.TotalIndices, sett.MemSettings.TotalTriangles, sett.MemSettings.TotalFaces, utilities.GetUsage()))

	imgui.EndMainMenuBar()
}
