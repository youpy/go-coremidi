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
  int n;
  char readbuffer[1024];

  while((n = read(fd, readbuffer, sizeof(readbuffer) - 1)) > 0) {
    readbuffer[n] = 0x00;

    goCallback(proc, source, readbuffer);
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
  int i, n;

  for (i = 0; i < packetCount; i++) {
    n = write(*(int *)srcConnRefCon, &packet->data, packet->length);
    packet = MIDIPacketNext(packet);
  }
}

typedef void (*midi_input_proc)(const MIDIPacketList *pktlist, void *readProcRefCon, void *srcConnRefCon);

static midi_input_proc getProc()
{
  return *MIDIInputProc;
}

*/
import "C"

//export goCallback
func goCallback(proc unsafe.Pointer, source unsafe.Pointer, p1 *C.char) {
	foo := *(*func(source Source, value []byte))(proc)

	foo(*(*Source)(source), ([]byte)(C.GoString(p1)))
}

type OutputPort struct {
	port C.MIDIPortRef
}

type InputPort struct {
	port     C.MIDIPortRef
	readProc func(source Source, value []byte)
}

func NewOutputPort(client Client, name string) (outputPort OutputPort, err error) {
	var port C.MIDIPortRef

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	osStatus := C.MIDIOutputPortCreate(client.client, C.CFStringCreateWithCString(nil, cName, C.kCFStringEncodingMacRoman), &port)

	if osStatus != C.noErr {
		err = errors.New(fmt.Sprintf("%d: failed to create a port", int(osStatus)))
	} else {
		outputPort = OutputPort{port}
	}

	return
}

func NewInputPort(client Client, name string, readProc func(source Source, value []byte)) (inputPort InputPort, err error) {
	var port C.MIDIPortRef

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	osStatus := C.MIDIInputPortCreate(client.client,
		C.CFStringCreateWithCString(nil, cName, C.kCFStringEncodingMacRoman),
		(C.MIDIReadProc)(unsafe.Pointer(C.getProc())),
		unsafe.Pointer(uintptr(0)),
		&port)

	if osStatus != C.noErr {
		err = errors.New(fmt.Sprintf("%d: failed to create a port", int(osStatus)))
	} else {
		inputPort = InputPort{port, readProc}
	}

	return
}

func (port InputPort) Connect(source Source) {
	fd := pipe()
	writeFd := C.getIntValue(fd, 1)

	C.MIDIPortConnectSource(port.port, source.endpoint, unsafe.Pointer(&writeFd))

	go func() {
		// TODO: should terminate when MIDIPortDisconnectSource is called
		C.readFromPipeAndCallback(C.getIntValue(fd, 0), unsafe.Pointer(&port.readProc), unsafe.Pointer(&source))
	}()
}

func pipe() *C.int {
	return C.make_pipe()
}
