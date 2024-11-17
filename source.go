package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
*/
import "C"
import (
	"errors"
	"fmt"
)

type Source struct {
	endpoint C.MIDIEndpointRef
	*Object
}

func NewSource(client Client, name string) (source Source, err error) {
	var endpointRef C.MIDIEndpointRef

	stringToCFString(name, func(cfName C.CFStringRef) {
		osStatus := C.MIDISourceCreate(client.client, cfName, &endpointRef)

		if osStatus != C.noErr {
			err = fmt.Errorf("%d: failed to create a source", int(osStatus))
		} else {
			source = Source{endpointRef, &Object{C.MIDIObjectRef(endpointRef)}}
		}
	})

	return
}

func AllSources() (sources []Source, err error) {
	numberOfSources := numberOfSources()
	sources = make([]Source, numberOfSources)

	for i := range sources {
		source := C.MIDIGetSource(C.ItemCount(i))

		if source == (C.MIDIEndpointRef)(0) {
			err = errors.New("failed to get source")

			return
		}

		sources[i] = Source{
			source,
			&Object{C.MIDIObjectRef(source)}}
	}

	return
}

func (source *Source) Entity() (entity Entity) {
	var entityRef C.MIDIEntityRef

	C.MIDIEndpointGetEntity(source.endpoint, &entityRef)
	entity = Entity{entityRef, &Object{C.MIDIObjectRef(entityRef)}}

	return
}

func numberOfSources() int {
	return int(C.ItemCount(C.MIDIGetNumberOfSources()))
}
