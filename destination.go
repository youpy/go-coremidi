package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
*/
import "C"
import "errors"

type Destination struct {
	endpoint C.MIDIEndpointRef
	*Object
}

func AllDestinations() (destinations []Destination, err error) {
	numberOfDestinations := numberOfDestinations()
	destinations = make([]Destination, numberOfDestinations)

	for i := range destinations {
		destination := C.MIDIGetDestination(C.ItemCount(i))

		if destination == (C.MIDIEndpointRef)(0) {
			err = errors.New("failed to get destination")

			return
		}

		destinations[i] = Destination{
			destination,
			&Object{C.MIDIObjectRef(destination)}}
	}

	return
}

func numberOfDestinations() int {
	return int(C.ItemCount(C.MIDIGetNumberOfDestinations()))
}
