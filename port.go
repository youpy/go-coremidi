package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
*/
import "C"
import "unsafe"
import "errors"

type Port struct {
	port C.MIDIPortRef
}

type OutputPort Port

func NewOutputPort(client Client, name string) (outputPort OutputPort, err error) {
	var port C.MIDIPortRef

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	osStatus := C.MIDIOutputPortCreate(client.client, C.CFStringCreateWithCString(nil, cName, C.kCFStringEncodingMacRoman), &port)

	if osStatus != C.noErr {
		err = errors.New("failed to create a port")
	} else {
		outputPort = OutputPort{port}
	}

	return
}
