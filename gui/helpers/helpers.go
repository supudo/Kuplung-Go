package helpers

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/gui/fonts"
)

// AddControlsSlider ...
func AddControlsSlider(title string, idx int32, step float32, min float32, limit float32, showAnimate bool, animatedFlag *bool, animatedValue *float32, doMinus bool, isFrame *bool) bool {
	if len(title) > 0 {
		imgui.Text(title)
	}
	if showAnimate {
		cid := fmt.Sprintf("##%v", idx)
		_ = imgui.Checkbox(cid, animatedFlag)
		if *animatedFlag {
			animateValue(isFrame, animatedFlag, animatedValue, step, limit, doMinus)
		}
		if imgui.IsItemHovered() {
			imgui.SetTooltip("Animate " + title)
		}
		imgui.SameLine()
	}
	sid := fmt.Sprintf("##%v", idx)
	return imgui.SliderFloat(sid, animatedValue, min, limit)
}

// AddControlColor3 ...
func AddControlColor3(title string, vValue *mgl32.Vec3, bValue *bool) {
	ceid := fmt.Sprintf("##101%v", title)
	iconid := fmt.Sprintf("%c %v", fonts.MD_ICON_COLORIZE, ceid)
	imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: vValue.X(), Y: vValue.Y(), Z: vValue.Z(), W: 1})
	imgui.Text(title)
	imgui.PopStyleColor()
	vv := [3]float32{vValue.X(), vValue.Y(), vValue.Z()}
	imgui.ColorEdit3(ceid, &vv)
	vValue[0] = vv[0]
	vValue[1] = vv[1]
	vValue[2] = vv[2]
	imgui.SameLine()
	imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 0})
	imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0})
	imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0})
	imgui.PushStyleColor(imgui.StyleColorBorder, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0})
	if imgui.ButtonV(iconid, imgui.Vec2{X: 0, Y: 0}) {
		*bValue = !*bValue
	}
	imgui.PopStyleColorV(4)
	if *bValue {
		imgui.ColorPicker3(title, &vv, 0)
	}
}

// AddControlColor4 ...
func AddControlColor4(title string, vValue *mgl32.Vec4, bValue *bool) {
	ceid := fmt.Sprintf("##101%v", title)
	iconid := fmt.Sprintf("%c %v", fonts.MD_ICON_COLORIZE, ceid)
	imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: vValue.X(), Y: vValue.Y(), Z: vValue.Z(), W: 1.0})
	imgui.Text(title)
	imgui.PopStyleColorV(1)
	vv := [4]float32{vValue.X(), vValue.Y(), vValue.Z(), vValue.W()}
	imgui.ColorEdit4(ceid, &vv)
	vValue[0] = vv[0]
	vValue[1] = vv[1]
	vValue[2] = vv[2]
	vValue[3] = vv[3]
	imgui.SameLine()
	imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0})
	imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0})
	imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0})
	imgui.PushStyleColor(imgui.StyleColorBorder, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0})
	if imgui.ButtonV(iconid, imgui.Vec2{X: 0, Y: 0}) {
		*bValue = !*bValue
	}
	imgui.PopStyleColorV(4)
	if *bValue {
		imgui.ColorPicker4(title, &vv, 0)
	}
}

// AddControlsSliderSameLine ...
func AddControlsSliderSameLine(title string, idx int32, step float32, min float32, limit float32, showAnimate bool, animatedFlag *bool, animatedValue *float32, doMinus bool, isFrame *bool) {
	if showAnimate {
		cid := fmt.Sprintf("##%v", idx)
		_ = imgui.Checkbox(cid, animatedFlag)
		animateValue(isFrame, animatedFlag, animatedValue, step, limit, doMinus)
		if imgui.IsItemHovered() {
			imgui.SetTooltip("Animate " + title)
		}
		imgui.SameLine()
	}
	sid := fmt.Sprintf("##10%v", idx)
	imgui.SliderFloat(sid, *(&animatedValue), min, limit)
	imgui.SameLine()
	imgui.Text(title)
}

// AddControlsIntegerSlider ...
func AddControlsIntegerSlider(title string, idx, min, limit int32, animatedValue *int32) {
	if len(title) > 0 {
		imgui.Text(fmt.Sprintf("%s", title))
	}
	sid := "##10" + fmt.Sprint(idx)
	imgui.SliderInt(sid, *(&animatedValue), min, limit)
}

// AddControlsFloatSlider ...
func AddControlsFloatSlider(title string, idx int32, min, limit float32, animatedValue *float32) {
	if len(title) > 0 {
		imgui.Text(title)
	}
	sid := "##10" + fmt.Sprint(idx)
	imgui.SliderFloat(sid, *(&animatedValue), min, limit)
}

func animateValue(isFrame, animatedFlag *bool, animatedValue *float32, step, limit float32, doMinus bool) {
	go animateValueAsync(isFrame, animatedFlag, animatedValue, step, limit, doMinus)
}

func animateValueAsync(isFrame, animatedFlag *bool, animatedValue *float32, step, limit float32, doMinus bool) {
	for *animatedFlag {
		if *isFrame {
			v := *animatedValue
			v += step
			if v > limit {
				if doMinus {
					v = -1 * limit
				} else {
					v = 0
				}
			}
			*animatedValue = v
			// TODO: fix the proper framerate for animated values
			// *isFrame = false
			time.Sleep(2500 * time.Millisecond)
		}
	}
}
