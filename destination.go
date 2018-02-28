package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
#include <stdio.h>
#include <unistd.h>

static void MIDIDestinationInputProc(const MIDIPacketList *pktlist, void *readProcRefCon, void *srcConnRefCon)
{
  MIDIPacket *packet = (MIDIPacket *)&(pktlist->packet[0]);
  UInt32 packetCount = pktlist->numPackets;
  int i, j, n;
  Byte *data;

  int lengthBytes = 2;
  int timeStampBytes = 8;

  for (i = 0; i < packetCount; i++) {
    data = calloc(sizeof(Byte), packet->length + lengthBytes + timeStampBytes);

    memcpy(data, &(packet->length), lengthBytes);
    memcpy(data + lengthBytes, &(packet->timeStamp), timeStampBytes);
    memcpy(data + lengthBytes + timeStampBytes, packet->data, packet->length);

    // http://man7.org/linux/man-pages/man7/pipe.7.html
    //
    // POSIX.1-2001 says that write(2)s of less than PIPE_BUF bytes must be
    // atomic: the output data is written to the pipe as a contiguous sequence.
    //
    // POSIX.1-2001 requires PIPE_BUF to be at least 512 bytes.
    n = write(*(int *)readProcRefCon, data, packet->length + lengthBytes + timeStampBytes);
    packet = MIDIPacketNext(packet);
    free(data);
  }
}

typedef void (*midi_destination_input_proc)(const MIDIPacketList *pktlist, void *readProcRefCon, void *srcConnRefCon);

static midi_destination_input_proc getMidiDestinationProc()
{
  return *MIDIDestinationInputProc;
}

*/
import "C"
import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

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

func NewDestination(client Client, name string, readProc func(packet Packet)) (destination Destination, err error) {
	var endpointRef C.MIDIEndpointRef

	fd := make([]int, 2)
	syscall.Pipe(fd)
	readFd := fd[0]
	writeFd := C.int(fd[1])

	go processImcomingPacket(
		readFd,
		func(data []byte, timeStamp uint64) {
			readProc(NewPacket(data, timeStamp))
		},
	)

	stringToCFString(name, func(cfName C.CFStringRef) {
		osStatus := C.MIDIDestinationCreate(
			client.client,
			cfName,
			(C.MIDIReadProc)(C.getMidiDestinationProc()),
			unsafe.Pointer(&writeFd),
			&endpointRef,
		)

		if osStatus != C.noErr {
			err = errors.New(fmt.Sprintf("%d: failed to create a destination", int(osStatus)))
		} else {
			destination = Destination{endpointRef, &Object{C.MIDIObjectRef(endpointRef)}}
		}
	})

	return
}

func (dest Destination) Dispose() {
	C.MIDIEndpointDispose(dest.endpoint)
}

func numberOfDestinations() int {
	return int(C.ItemCount(C.MIDIGetNumberOfDestinations()))
}
