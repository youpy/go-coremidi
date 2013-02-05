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

type Packet struct {
	packetList C.MIDIPacketList
}

func NewPacket(values ...int) Packet {
	var packetList C.MIDIPacketList
	var data = C.makeByteArray(C.int(len(values)))
	defer C.freeByteArray(data)

	for i := range values {
		C.setByte(data, (C.Byte)(values[i]), C.int(i))
	}

	packet := C.MIDIPacketListInit(&packetList)
	packet = C.MIDIPacketListAdd(&packetList, 1024, packet, 0, C.ByteCount(len(values)), data)

	return Packet{packetList}
}

func (packet Packet) Send(port OutputPort, destination Destination) (result int, err error) {
	osStatus := C.MIDISend(port.port, destination.endpoint, &packet.packetList)

	if osStatus != C.noErr {
		err = errors.New("failed to send MIDI")
	} else {
		result = int(osStatus)
	}

	return
}
