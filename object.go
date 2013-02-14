package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
*/
import "C"

type Object struct {
	object C.MIDIObjectRef
}

func (object Object) Name() string {
	return object.getStringProperty(C.kMIDIPropertyName)
}

func (object Object) Manufacturer() string {
	return object.getStringProperty(C.kMIDIPropertyManufacturer)
}

func (object Object) getStringProperty(key C.CFStringRef) (propValue string) {
	var result C.CFStringRef

	osStatus := C.MIDIObjectGetStringProperty(object.object, key, &result)

	if osStatus != C.noErr {
		return
	}

	defer C.CFRelease((C.CFTypeRef)(result))

	value := C.CFStringGetCStringPtr(result, C.kCFStringEncodingMacRoman)
	propValue = C.GoString(value)

	return
}
