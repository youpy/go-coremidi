package coremidi

import "unsafe"
import "errors"
import "fmt"

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
#include <stdio.h>
#include <unistd.h>

extern void goCallback(void *proc, void *source, char *value);

static int getIntValue(int *array, int index)
{
  return array[index];
}

static void readFromPipeAndCallback(int fd, void *proc, void *source)
{
  int n, size;
  char readbuffer[30];

  while((n = read(fd, readbuffer, 1)) > 0) {
    size = readbuffer[0];
    n = read(fd, readbuffer, size);

    if(n == size) {
      readbuffer[size] = 0x00;
      goCallback(proc, source, readbuffer);
    }
  }
}

static int *make_pipe()
{
  int *fd = calloc(sizeof(int), 2);

  pipe(fd);

  return fd;
}

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
	writeFds []*C.int
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
			inputPort = InputPort{port, readProc, make([]*C.int, 0)}
		}
	})

	return
}

func (port InputPort) Connect(source Source) (portConnection, error) {
	fd := pipe()
	readFd := C.getIntValue(fd, 0)
	writeFd := C.getIntValue(fd, 1)
	port.writeFds = append(port.writeFds, &writeFd)

	C.MIDIPortConnectSource(port.port, source.endpoint, unsafe.Pointer(&writeFd))

	go func() {
		C.readFromPipeAndCallback(
			readFd,
			unsafe.Pointer(&port.readProc),
			unsafe.Pointer(&source))

		C.close(readFd)
	}()

	return portConnection{port, source, &writeFd}, nil
}

type portConnection struct {
	port    InputPort
	source  Source
	writeFd *C.int
}

func (connection portConnection) Disconnect() {
	C.close(*connection.writeFd)
	C.MIDIPortDisconnectSource(connection.port.port, connection.source.endpoint)
}

func pipe() *C.int {
	return C.make_pipe()
}

//export goCallback
func goCallback(proc unsafe.Pointer, source unsafe.Pointer, p1 *C.char) {
	foo := *(*ReadProc)(proc)

	foo(*(*Source)(source), ([]byte)(C.GoString(p1)))
}
