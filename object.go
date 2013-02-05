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

func (object Object) getStringProperty(key C.CFStringRef) string {
	var result C.CFStringRef

	C.MIDIObjectGetStringProperty(object.object, key, &result)
	value := C.CFStringGetCStringPtr(result, C.kCFStringEncodingMacRoman)

	return C.GoString(value)
}
