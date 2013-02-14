package coremidi

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"
import "unsafe"

func stringToCFString(str string, callback func(cfStr C.CFStringRef)) {
	cStr := C.CString(str)
	cfStr := C.CFStringCreateWithCString(nil, cStr, C.kCFStringEncodingMacRoman)

	defer C.free(unsafe.Pointer(cStr))
	defer C.CFRelease((C.CFTypeRef)(cfStr))

	callback(cfStr)
}
