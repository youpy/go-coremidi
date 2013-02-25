package coremidi

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
#include <stdio.h>
#include <unistd.h>

static void MIDIInputProc(const MIDIPacketList *pktlist, void *readProcRefCon,  void *srcConnRefCon)
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

    n = write(*(int *)srcConnRefCon, data, packet->length + 1);
    packet = MIDIPacketNext(packet);
    free(data);
  }
}

typedef void (*midi_input_proc)(const MIDIPacketList *pktlist, void *readProcRefCon, void *srcConnRefCon);

static midi_input_proc getProc()
{
  return *MIDIInputProc;
}

*/
import "C"

type OutputPort struct {
	port C.MIDIPortRef
}

func NewOutputPort(client Client, name string) (outputPort OutputPort, err error) {
	var port C.MIDIPortRef

	stringToCFString(name, func(cfName C.CFStringRef) {
		osStatus := C.MIDIOutputPortCreate(client.client, cfName, &port)

		if osStatus != C.noErr {
			err = errors.New(fmt.Sprintf("%d: failed to create a port", int(osStatus)))
		} else {
			outputPort = OutputPort{port}
		}
	})

	return
}

type ReadProc func(source Source, value []byte)
type InputPort struct {
	port     C.MIDIPortRef
	readProc ReadProc
}

func NewInputPort(client Client, name string, readProc ReadProc) (inputPort InputPort, err error) {
	var port C.MIDIPortRef

	stringToCFString(name, func(cfName C.CFStringRef) {
		osStatus := C.MIDIInputPortCreate(client.client,
			cfName,
			(C.MIDIReadProc)(C.getProc()),
			unsafe.Pointer(uintptr(0)),
			&port)

		if osStatus != C.noErr {
			err = errors.New(fmt.Sprintf("%d: failed to create a port", int(osStatus)))
		} else {
			inputPort = InputPort{port, readProc}
		}
	})

	return
}

func (port InputPort) Connect(source Source) (portConnection, error) {
	fd := make([]int, 2)

	syscall.Pipe(fd)

	readFd := fd[0]
	writeFd := C.int(fd[1])

	C.MIDIPortConnectSource(port.port, source.endpoint, unsafe.Pointer(&writeFd))

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

			port.readProc(source, data)
		}

		syscall.Close(readFd)
	}()

	return portConnection{port, source, &writeFd}, nil
}

type portConnection struct {
	port    InputPort
	source  Source
	writeFd *C.int
}

func (connection portConnection) Disconnect() {
	syscall.Close(int(*connection.writeFd))
	C.MIDIPortDisconnectSource(connection.port.port, connection.source.endpoint)
}
