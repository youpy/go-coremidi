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

  for (i = 0; i < packetCount; i++) {
    data = calloc(sizeof(Byte), packet->length + 1);
    *data = packet->length;

    for (j = 0; j < packet->length; j++) {
      *(data + j + 1) = *(packet->data + j);
    }

    // http://man7.org/linux/man-pages/man7/pipe.7.html
    //
    // POSIX.1-2001 says that write(2)s of less than PIPE_BUF bytes must be
    // atomic: the output data is written to the pipe as a contiguous sequence.
    //
    // POSIX.1-2001 requires PIPE_BUF to be at least 512 bytes.
    n = write(*(int *)readProcRefCon, data, packet->length + 1);
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

func NewDestination(client Client, name string, readProc func(value []byte)) (destination Destination, err error) {
	var endpointRef C.MIDIEndpointRef

	fd := make([]int, 2)
	syscall.Pipe(fd)
	readFd := fd[0]
	writeFd := C.int(fd[1])

	go func() {
		dataForLength := make([]byte, 1)

		for {
			n, err := syscall.Read(readFd, dataForLength)
			if err != nil || n != 1 {
				break
			}

			length := dataForLength[0]
			data := make([]byte, length)

			n, err = syscall.Read(readFd, data)
			if err != nil || n != int(length) {
				break
			}

			readProc(data)
		}

		syscall.Close(readFd)
	}()

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
