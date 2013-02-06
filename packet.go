package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>

Byte *makeByteArray(int size) {
  return calloc(sizeof(Byte), size);
}

void setByte(Byte *array, Byte value, int n) {
  array[n] = value;
}

void freeByteArray(Byte *array) {
 free(array);
}
*/
import "C"
import "errors"
import "fmt"

type Packet struct {
	packetList C.MIDIPacketList
}

func NewPacket(p []byte) Packet {
	var packetList C.MIDIPacketList
	var data = C.makeByteArray(C.int(len(p)))
	defer C.freeByteArray(data)

	for i := range p {
		C.setByte(data, (C.Byte)(p[i]), C.int(i))
	}

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
