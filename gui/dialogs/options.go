package dialogs

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/settings"
)

// ViewOptions ...
type ViewOptions struct {
}

// NewViewOptions ...
func NewViewOptions() *ViewOptions {
	return &ViewOptions{}
}

// Render ...
func (view *ViewOptions) Render(open, isFrame *bool) {
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()

	imgui.SetNextWindowSizeV(imgui.Vec2{X: 400, Y: 560}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: 200, Y: 200}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	if imgui.BeginV("Options", open, imgui.WindowFlagsResizeFromAnySide) {
		if imgui.TreeNodeV("General", imgui.TreeNodeFlagsCollapsingHeader) {
			imgui.TreePop()
		}
		if imgui.TreeNodeV("Rendering", imgui.TreeNodeFlagsCollapsingHeader) {
			if imgui.Checkbox("Show OpenGL Errors", &sett.Rendering.ShowGLErrors) {
				settings.SaveSettings()
			}
			if imgui.Checkbox("Occlusion Culling", &rsett.General.OcclusionCulling) {
				settings.SaveRenderingSettings()
			}
			availableRenderers := []string{"Simple", "Forward", "Forward with Shadow Mapping", "Shadow Mapping", "Deferred"}
			if imgui.BeginCombo("Renderer", availableRenderers[sett.App.RendererType]) {
				var i uint32
				for i = 0; i < uint32(len(availableRenderers)); i++ {
					rsel := (i == sett.App.RendererType)
					if imgui.SelectableV(availableRenderers[i], rsel, 0, imgui.Vec2{0, 0}) {
						sett.App.RendererType = i
						settings.SaveSettings()
					}
					if rsel {
						imgui.SetItemDefaultFocus()
					}
				}
				imgui.EndCombo()
			}
			imgui.TreePop()
		}
		if imgui.TreeNodeV("Look & Feel", imgui.TreeNodeFlagsCollapsingHeader) {
			imgui.TreePop()
		}
		imgui.End()
	}
}
