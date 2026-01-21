package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation -framework CoreServices
#include <CoreMIDI/CoreMIDI.h>
#include <CoreServices/CoreServices.h>
#include <stdint.h>

extern void midiNotifyCallback(uintptr_t handle, MIDINotification *message);
extern void midiNotifyProc(const MIDINotification *message, void *refCon);
*/
import "C"

import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

type Client struct {
	client       C.MIDIClientRef
	notifyHandle cgo.Handle
	hasNotify    bool
}

type MIDINotificationMessageID int32

const (
	MIDIMsgSetupChanged         MIDINotificationMessageID = 1
	MIDIMsgObjectAdded          MIDINotificationMessageID = 2
	MIDIMsgObjectRemoved        MIDINotificationMessageID = 3
	MIDIMsgPropertyChanged      MIDINotificationMessageID = 4
	MIDIMsgThruConnectionsChanged MIDINotificationMessageID = 5
	MIDIMsgSerialPortOwnerChanged MIDINotificationMessageID = 6
	MIDIMsgIOError              MIDINotificationMessageID = 7
)

func (id MIDINotificationMessageID) String() string {
	switch id {
	case MIDIMsgSetupChanged:
		return "SetupChanged"
	case MIDIMsgObjectAdded:
		return "ObjectAdded"
	case MIDIMsgObjectRemoved:
		return "ObjectRemoved"
	case MIDIMsgPropertyChanged:
		return "PropertyChanged"
	case MIDIMsgThruConnectionsChanged:
		return "ThruConnectionsChanged"
	case MIDIMsgSerialPortOwnerChanged:
		return "SerialPortOwnerChanged"
	case MIDIMsgIOError:
		return "IOError"
	default:
		return "Unknown"
	}
}

type Notification struct {
	MessageID MIDINotificationMessageID
}

// TODO: Expose structured notification payloads (e.g. add/remove/property change).
type NotifyFunc func(notification Notification)

// NewClient creates a CoreMIDI client without a notify callback.
// A CFRunLoop is not required to call polling APIs, but without a running loop
// on the client's thread, CoreMIDI may not deliver device list updates, so
// polling can return stale results in long-running processes.
func NewClient(name string) (client Client, err error) {
	var clientRef C.MIDIClientRef

	stringToCFString(name, func(cfName C.CFStringRef) {
		osStatus := C.MIDIClientCreate(cfName, nil, nil, &clientRef)

		if osStatus != C.noErr {
			err = fmt.Errorf("%d: failed to create a client", int(osStatus))
		} else {
			client = Client{client: clientRef}
		}
	})

	return
}

// NewClientWithNotification creates a CoreMIDI client and registers a notify callback.
// CoreMIDI delivers device updates via a CFRunLoop on the same OS thread that created
// the client. That loop must be running for the system's device list to update at all
// (including for polling APIs), not just for notifications to fire.
// See: https://stackoverflow.com/questions/19002598/midiclientcreate-notifyproc-not-getting-called
func NewClientWithNotification(name string, notify NotifyFunc) (client Client, err error) {
	var clientRef C.MIDIClientRef
	handle := cgo.NewHandle(notify)
	refCon := unsafe.Pointer(handle)

	stringToCFString(name, func(cfName C.CFStringRef) {
		osStatus := C.MIDIClientCreate(cfName, (C.MIDINotifyProc)(C.midiNotifyProc), refCon, &clientRef)

		if osStatus != C.noErr {
			err = fmt.Errorf("%d: failed to create a client", int(osStatus))
		} else {
			client = Client{client: clientRef, notifyHandle: handle, hasNotify: true}
		}
	})

	if err != nil {
		handle.Delete()
	}

	return
}

func (client *Client) Close() error {
	if client.client == 0 {
		return nil
	}

	osStatus := C.MIDIClientDispose(client.client)
	client.client = 0

	if client.hasNotify {
		client.notifyHandle.Delete()
		client.hasNotify = false
	}

	if osStatus != C.noErr {
		return fmt.Errorf("%d: failed to dispose client", int(osStatus))
	}

	return nil
}

//export midiNotifyCallback
func midiNotifyCallback(handle C.uintptr_t, message *C.MIDINotification) {
	var messageID MIDINotificationMessageID
	if message != nil {
		messageID = MIDINotificationMessageID(message.messageID)
	}
	if handle == 0 {
		return
	}

	h := cgo.Handle(handle)
	if notify, ok := h.Value().(NotifyFunc); ok && notify != nil {
		notify(Notification{
			MessageID: messageID,
		})
	}
}
