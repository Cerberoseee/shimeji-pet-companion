//go:build windows

package physics

import "syscall"

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

const vkLButton = 0x01

func isLButtonPressed() bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vkLButton))
	return (ret & 0x8000) != 0
}
