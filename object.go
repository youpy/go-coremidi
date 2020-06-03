package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>

// based on https://stackoverflow.com/a/9166500
char * MyCFStringCopyUTF8String(CFStringRef aString, CFStringEncoding encoding) {
  if (aString == NULL) {
    return NULL;
  }

  CFIndex length = CFStringGetLength(aString);
  CFIndex maxSize =
    CFStringGetMaximumSizeForEncoding(length, encoding) + 1;
  char *buffer = (char *)malloc(maxSize);

  if (CFStringGetCString(aString, buffer, maxSize,
                         encoding)) {
    return buffer;
  }

  // If we failed
  free(buffer);
  return NULL;
}

*/
import "C"
import "unsafe"

type Object struct {
	object C.MIDIObjectRef
}

func (object Object) Name() string {
	return object.getStringProperty(C.kMIDIPropertyName)
}

func (object Object) Manufacturer() string {
	return object.getStringProperty(C.kMIDIPropertyManufacturer)
}

func (object Object) getStringProperty(key C.CFStringRef) (propValue string) {
	var result C.CFStringRef

	osStatus := C.MIDIObjectGetStringProperty(object.object, key, &result)

	if osStatus != C.noErr {
		return
	}

	defer C.CFRelease((C.CFTypeRef)(result))

	value := C.MyCFStringCopyUTF8String(result, C.kCFStringEncodingUTF8)
	defer C.free(unsafe.Pointer(value))

	propValue = C.GoString(value)

	return
}
