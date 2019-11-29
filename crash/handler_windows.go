package crash

import "C"

import (
	"fmt"
	"runtime/debug"
	"unsafe"
)

func init() {
	Handler = windowsHandler
}

func windowsHandler(versionInfo string) {
	if e := recover(); e != nil {
		text := `Kuplung has crashed.
		Details follow bellow:
`
		text += fmt.Sprintf("%s:\n%s", e, debug.Stack())
		messageBox(text, "Something unexpected happened - "+versionInfo)
	}
}

func messageBox(text, caption string) {
	textArg, textFin := wrapString(text)
	defer textFin()
	captionArg, captionFin := wrapString(caption)
	defer captionFin()
	var hwnd unsafe.Pointer
	C.MessageBoxExA((*C.struct_HWND__)(hwnd), textArg, captionArg, C.UINT(0), C.WORD(0))
}

func wrapString(value string) (wrapped *C.char, finisher func()) {
	wrapped = C.CString(value)
	finisher = func() { C.free(unsafe.Pointer(wrapped)) } // nolint: gas
	return
}
