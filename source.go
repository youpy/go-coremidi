package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
*/
import "C"
import "unsafe"
import "errors"
import "fmt"

type Source struct {
	endpoint C.MIDIEndpointRef
	*Object
}

func NewSource(client Client, name string) (source Source, err error) {
	var endpointRef C.MIDIEndpointRef

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	osStatus := C.MIDISourceCreate(client.client, C.CFStringCreateWithCString(nil, cName, C.kCFStringEncodingMacRoman), &endpointRef)

	if osStatus != C.noErr {
		err = errors.New(fmt.Sprintf("%d: failed to create a source", int(osStatus)))
	} else {
		source = Source{endpointRef, &Object{C.MIDIObjectRef(endpointRef)}}
	}

	return
}
