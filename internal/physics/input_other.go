//go:build !windows

package physics

/*
#cgo LDFLAGS: -framework CoreGraphics
#include <CoreGraphics/CoreGraphics.h>
*/
import "C"

// isLButtonPressed returns false on platforms where native mouse-button polling
// is not implemented. Custom drag events continue to work cross-platform.
func isLButtonPressed() bool {
	state := C.CGEventSourceButtonState(
		C.CGEventSourceStateID(0),
		C.kCGMouseButtonLeft,
	)
	return bool(state)
}
