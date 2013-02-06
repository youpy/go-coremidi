package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
*/
import "C"
import "errors"

type Entity struct {
	entity C.MIDIEntityRef
	*Object
}

func (entity Entity) Sources() (sources []Source, err error) {
	numberOfSources := int(C.ItemCount(C.MIDIEntityGetNumberOfSources(entity.entity)))
	sources = make([]Source, numberOfSources)

	for i := range sources {
		source := C.MIDIEntityGetSource(entity.entity, C.ItemCount(i))

		if source == (C.MIDIEndpointRef)(0) {
			err = errors.New("failed to get source")

			return
		}

		sources[i] = Source{source, &Object{C.MIDIObjectRef(source)}}
	}

	return
}

func (entity Entity) Destinations() (destinations []Destination, err error) {
	numberOfDestinations := int(C.ItemCount(C.MIDIEntityGetNumberOfDestinations(entity.entity)))
	destinations = make([]Destination, numberOfDestinations)

	for i := range destinations {
		destination := C.MIDIEntityGetDestination(entity.entity, C.ItemCount(i))

		if destination == (C.MIDIEndpointRef)(0) {
			err = errors.New("failed to get destination")

			return
		}

		destinations[i] = Destination{destination, &Object{C.MIDIObjectRef(destination)}}
	}

	return
}
