package coremidi

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"syscall"
	"unsafe"
)

func stringToCFString(str string, callback func(cfStr C.CFStringRef)) {
	cStr := C.CString(str)
	cfStr := C.CFStringCreateWithCString(C.kCFAllocatorDefault, cStr, C.kCFStringEncodingMacRoman)

	defer C.free(unsafe.Pointer(cStr))
	defer C.CFRelease((C.CFTypeRef)(cfStr))

	callback(cfStr)
}

func processImcomingPacket(readFd int, onMessage func(data []byte, timeStamp uint64)) {
	var length uint16
	var timeStamp uint64

	defer syscall.Close(readFd)

	for {
		lengthBytes := make([]byte, 2)
		timeStampBytes := make([]byte, 8)

		n, err := syscall.Read(readFd, lengthBytes)
		if err != nil || n != 2 {
			break
		}

		err = binary.Read(bytes.NewBuffer(lengthBytes[:]), binary.LittleEndian, &length)
		if err != nil {
			break
		}

		n, err = syscall.Read(readFd, timeStampBytes)
		if err != nil || n != 8 {
			break
		}

		err = binary.Read(bytes.NewBuffer(timeStampBytes[:]), binary.LittleEndian, &timeStamp)
		if err != nil {
			break
		}

		data := make([]byte, length)

		n, err = syscall.Read(readFd, data)
		if err != nil || n != int(length) {
			break
		}

		onMessage(data, timeStamp)
	}
}
