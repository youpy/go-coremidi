package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI
#include <CoreMIDI/CoreMIDI.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Packet struct {
	Data      []byte
	TimeStamp uint64
}

func NewPacket(data []byte, timeStamp uint64) Packet {
	return Packet{data, timeStamp}
}

func (packet *Packet) createPacketList() C.MIDIPacketList {
	var packetList C.MIDIPacketList
	var data = (*C.Byte)(unsafe.Pointer(&packet.Data[0]))

	p := C.MIDIPacketListInit(&packetList)
	C.MIDIPacketListAdd(&packetList, 1024, p, C.MIDITimeStamp(packet.TimeStamp), C.ByteCount(len(packet.Data)), data)

	return packetList
}

func (packet *Packet) Send(port *OutputPort, destination *Destination) (err error) {
	packetList := packet.createPacketList()
	osStatus := C.MIDISend(port.port, destination.endpoint, &packetList)

	if osStatus != C.noErr {
		err = fmt.Errorf("%d: failed to send MIDI", int(osStatus))
	}

	return
}

func (packet *Packet) Received(source *Source) (err error) {
	packetList := packet.createPacketList()
	osStatus := C.MIDIReceived(source.endpoint, &packetList)

	if osStatus != C.noErr {
		err = fmt.Errorf("%d: failed to transmit MIDI", int(osStatus))
	}

	return
}
