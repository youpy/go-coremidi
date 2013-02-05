package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
*/
import "C"
import "errors"

type Device struct {
	device C.MIDIDeviceRef
	*Object
}

func AllDevices() (devices []Device, err error) {
	numberOfDevices := numberOfDevices()
	devices = make([]Device, numberOfDevices)

	for i := 0; i < numberOfDevices; i++ {
		device := C.MIDIGetDevice(C.ItemCount(i))

		if device == (C.MIDIDeviceRef)(0) {
			err = errors.New("failed to get device")

			return
		}

		devices[i] = Device{
			device,
			&Object{C.MIDIObjectRef(device)}}
	}

	return
}

func (device Device) Entities() (entities []Entity, err error) {
	numberOfEntitiles := int(C.ItemCount(C.MIDIDeviceGetNumberOfEntities(device.device)))
	entities = make([]Entity, numberOfEntitiles)

	for i := 0; i < numberOfEntitiles; i++ {
		entity := C.MIDIDeviceGetEntity(device.device, C.ItemCount(i))

		if entity == (C.MIDIEntityRef)(0) {
			err = errors.New("failed to get entity")

			return
		}

		entities[i] = Entity{entity, &Object{C.MIDIObjectRef(entity)}}
	}

	return
}

func numberOfDevices() int {
	return int(C.ItemCount(C.MIDIGetNumberOfDevices()))
}
