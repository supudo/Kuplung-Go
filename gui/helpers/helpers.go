package helpers

import (
	"fmt"
	"time"

	"github.com/inkyblackness/imgui-go"
)

// AddSliderF32 ...
func AddSliderF32(title string, idx int, step, min, limit float32, showAnimate, doMinus bool, animatedFlag, isFrame *bool, animatedValue *float32) {
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
	imgui.SliderFloat(sid, animatedValue, min, limit)
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
