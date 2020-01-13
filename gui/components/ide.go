package components

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// ComponentIDE ...
type ComponentIDE struct {
	selectedItem    int32
	selectedContent string
	items           []string
}

// NewComponentIDE ...
func NewComponentIDE() *ComponentIDE {
	cide := &ComponentIDE{}
	cide.selectedItem = 0
	cide.items = []string{
		"-- SELECT ITEM ---",
		"General - Vertex Shader",
		"General - Geometry Shader",
		"General - Tessellation Control Shader",
		"General - Tessellation Evaluation Shader",
		"General - Fragment Shader"}
	return cide
}

// Render ...
func (comp *ComponentIDE) Render(open *bool) {
	sett := settings.GetSettings()

	imgui.SetNextWindowSizeV(imgui.Vec2{X: sett.AppWindow.LogWidth, Y: sett.AppWindow.LogHeight}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: 40, Y: 40}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	if imgui.BeginV("IDE", open, imgui.WindowFlagsResizeFromAnySide) {
		if imgui.BeginCombo("##shader_ide", comp.items[comp.selectedItem]) {
			var i int32
			for i = 0; i < int32(len(comp.items)); i++ {
				sksel := (i == comp.selectedItem)
				if imgui.SelectableV(comp.items[i], sksel, 0, imgui.Vec2{X: 0, Y: 0}) {
					comp.selectedItem = i - 1
				}
				if sksel {
					imgui.SetItemDefaultFocus()
				}
			}
			imgui.EndCombo()
		}
		imgui.SameLine()
		if imgui.Button("Load Shader") {
			comp.loadSelectedItem()
		}
		imgui.Separator()

		if imgui.ButtonV("Compile Shaders", imgui.Vec2{X: -1.0, Y: 40.0}) {
			if sett.App.RendererType == types.InAppRendererTypeForward {
				if comp.selectedItem == 0 {
					sett.Components.ShaderSourceVertex = comp.selectedContent
				} else if comp.selectedItem == 1 {
					sett.Components.ShaderSourceGeometry = comp.selectedContent
				} else if comp.selectedItem == 2 {
					sett.Components.ShaderSourceTCS = comp.selectedContent
				} else if comp.selectedItem == 3 {
					sett.Components.ShaderSourceTES = comp.selectedContent
				} else if comp.selectedItem == 4 {
					sett.Components.ShaderSourceFragment = comp.selectedContent
				}
				sett.Components.ShouldRecompileShaders = true
			}
		}
		imgui.Separator()

		ws := imgui.WindowSize()
		imgui.InputTextMultilineV("", &comp.selectedContent, imgui.Vec2{X: ws.X, Y: ws.Y - 110}, 0, nil)

		imgui.End()
	}
}

func (comp *ComponentIDE) loadSelectedItem() {
	sett := settings.GetSettings()

	if sett.App.RendererType == types.InAppRendererTypeForward {
		if comp.selectedItem == 0 {
			comp.selectedContent = sett.Components.ShaderSourceVertex
		} else if comp.selectedItem == 1 {
			comp.selectedContent = sett.Components.ShaderSourceGeometry
		} else if comp.selectedItem == 2 {
			comp.selectedContent = sett.Components.ShaderSourceTCS
		} else if comp.selectedItem == 3 {
			comp.selectedContent = sett.Components.ShaderSourceTES
		} else if comp.selectedItem == 4 {
			comp.selectedContent = sett.Components.ShaderSourceFragment
		}
	}
}
