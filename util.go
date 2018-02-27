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
	cfStr := C.CFStringCreateWithCString(nil, cStr, C.kCFStringEncodingMacRoman)

	defer C.free(unsafe.Pointer(cStr))
	defer C.CFRelease((C.CFTypeRef)(cfStr))

	callback(cfStr)
}

func processImcomingPacket(readFd int, onFinish func(data []byte, timeStamp uint64)) {
	var timeStamp uint64

	dataForLength := make([]byte, 1)

	defer syscall.Close(readFd)

	for {
		n, err := syscall.Read(readFd, dataForLength)
		if err != nil || n != 1 {
			break
		}

		length := dataForLength[0]
		timeStampBytes := make([]byte, 8)

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

		onFinish(data, timeStamp)
	}
}
