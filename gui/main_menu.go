package gui

import (
	"fmt"
	"os"
	"runtime"

	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
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
		// lbl := ""
		// lbl += fmt.Sprintf("%#U", fonts.FA_ICON_FILE_O) + " "
		// lbl += fmt.Sprintf("%q", '\uf001') + " "
		// lbl += fmt.Sprintf("%v", fonts.FA_ICON_FILE_O) + " "
		// lbl += "New"
		if imgui.MenuItem("New") { // fonts.FA_ICON_FILE_O
			_, _ = trigger.Fire(types.ActionGuiActionFileNew)
		}
		if imgui.MenuItem("Open ...") { // fonts.FA_ICON_FOLDER_OPEN_O
			context.GuiVars.showOpenDialog = true
		}
		if imgui.BeginMenu("Open Recent") { // fonts.FA_ICON_FILES_O
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

		if imgui.MenuItem("Save ...") { // fonts.FA_ICON_FLOPPY_O
			context.GuiVars.showSaveDialog = true
		}

		imgui.Separator()

		if imgui.BeginMenu("Import") {
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

		if imgui.BeginMenu("Import Recent") { // fonts.FA_ICON_FILES_O
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

		if imgui.BeginMenu("Export") {
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
		if imgui.MenuItemV("Quit", quitShortcut, false, true) {
			os.Exit(3)
		}
		imgui.EndMenu()
	}

	if imgui.BeginMenu("Scene") {
		if imgui.BeginMenu("Add Light") {
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
		if imgui.BeginMenu("Scene Rendering") {
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
		imgui.MenuItemV("Render Image", "", context.GuiVars.showImageSave, true)
		imgui.MenuItemV("Renderer UI", "", context.GuiVars.showRendererUI, true)

		imgui.EndMenu()
	}

	if imgui.BeginMenu("View") {
		if imgui.MenuItem("Models") {
			context.GuiVars.showModels = !context.GuiVars.showModels
		}
		if imgui.MenuItem("Controls") {
			context.GuiVars.showControls = !context.GuiVars.showControls
		}
		lblVisualArtefacts := ""
		if rsett.General.ShowAllVisualArtefacts {
			lblVisualArtefacts = "Hide Visual Artefacts"
		} else {
			lblVisualArtefacts = "Show Visual Artefacts"
		}
		if imgui.MenuItem(lblVisualArtefacts) {
			rsett.General.ShowAllVisualArtefacts = !rsett.General.ShowAllVisualArtefacts
		}
		imgui.Separator()
		if imgui.MenuItem("Show Log Window") { // fonts.FA_ICON_BUG
			context.GuiVars.showLog = !context.GuiVars.showLog
		}
		if sett.App.RendererType == types.InAppRendererTypeForward {
			if imgui.MenuItem("IDE") { // fonts.FA_ICON_PENCIL
				context.GuiVars.showKuplungIDE = !context.GuiVars.showKuplungIDE
			}
		}
		if imgui.MenuItem("Screenshot") { // fonts.FA_ICON_DESKTOP
			context.GuiVars.showScreenshotWindow = !context.GuiVars.showScreenshotWindow
		}
		if imgui.MenuItem("Scene Statistics") { // fonts.FA_ICON_TACHOMETER
			context.GuiVars.showSceneStats = !context.GuiVars.showSceneStats
		}
		if imgui.MenuItem("Structured Volumetric Sampling") { // fonts.FA_ICON_PAPER_PLANE_O
			context.GuiVars.showSVS = !context.GuiVars.showSVS
		}
		if imgui.MenuItem("Shadertoy") { // fonts.FА_ICON_BICYCLE
			context.GuiVars.showShadertoy = !context.GuiVars.showShadertoy
		}
		imgui.Separator()
		if imgui.MenuItem("Options") {
			context.GuiVars.showOptions = !context.GuiVars.showOptions
		}
		imgui.Separator()
		if imgui.MenuItem("Cube") {
			rsett.General.ShowCube = !rsett.General.ShowCube
		}
		imgui.EndMenu()
	}

	if imgui.BeginMenu("Help") {
		if imgui.MenuItem("Metrics") {
			context.GuiVars.showMetrics = !context.GuiVars.showMetrics
		}
		if imgui.MenuItem("About ImGui") {
			context.GuiVars.showAboutImGui = !context.GuiVars.showAboutImGui
		}
		if imgui.MenuItem("About Kuplung") {
			context.GuiVars.showAboutKuplung = !context.GuiVars.showAboutKuplung
		}
		imgui.Separator()
		if imgui.MenuItem("ImGui Demo Window") {
			context.GuiVars.showDemoWindow = !context.GuiVars.showDemoWindow
		}
		imgui.EndMenu()
	}

	imgui.Text(fmt.Sprintf("  | [%.4f ms/frame] %d objs, %d verts, %d indices (%d tris, %d faces) | %v", sett.MemSettings.NbResult, sett.MemSettings.TotalObjects, sett.MemSettings.TotalVertices, sett.MemSettings.TotalIndices, sett.MemSettings.TotalTriangles, sett.MemSettings.TotalFaces, utilities.GetUsage()))

	imgui.EndMainMenuBar()
}
