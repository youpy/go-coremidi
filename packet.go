package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI
#include <CoreMIDI/CoreMIDI.h>
*/
import "C"
import "errors"
import "fmt"
import "unsafe"

type Packet struct {
	packetList C.MIDIPacketList
}

func NewPacket(p []byte) Packet {
	var packetList C.MIDIPacketList
	var data = (*C.Byte)(unsafe.Pointer(&p[0]))

	packet := C.MIDIPacketListInit(&packetList)
	packet = C.MIDIPacketListAdd(&packetList, 1024, packet, 0, C.ByteCount(len(p)), data)

	return Packet{packetList}
}

func (packet Packet) Send(port *OutputPort, destination *Destination) (err error) {
	osStatus := C.MIDISend(port.port, destination.endpoint, &packet.packetList)

	if osStatus != C.noErr {
		err = errors.New(fmt.Sprintf("%d: failed to send MIDI", int(osStatus)))
	}

	return
}

func (packet Packet) Received(source *Source) (err error) {
	osStatus := C.MIDIReceived(source.endpoint, &packet.packetList)

	if osStatus != C.noErr {
		err = errors.New(fmt.Sprintf("%d: failed to transmit MIDI", int(osStatus)))
	}

	return
}
