package coremidi

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>

#include <string.h>

CFRunLoopRef currentRunLoop() {
  return CFRunLoopGetCurrent();
}

void retainRunLoop(CFRunLoopRef loop) {
  if (loop != NULL) {
    CFRetain(loop);
  }
}

void releaseRunLoop(CFRunLoopRef loop) {
  if (loop != NULL) {
    CFRelease(loop);
  }
}

void runLoopRun() {
  CFRunLoopRun();
}

void runLoopStop(CFRunLoopRef loop) {
  if (loop != NULL) {
    CFRunLoopStop(loop);
    CFRunLoopWakeUp(loop);
  }
}

static void stopSourcePerform(void *info) {
	(void)info;
	CFRunLoopStop(CFRunLoopGetCurrent());
}

CFRunLoopSourceRef createStopSource(CFRunLoopRef loop) {
	if (loop == NULL) {
		return NULL;
	}
	CFRunLoopSourceContext ctx = {0};
	ctx.perform = stopSourcePerform;
	CFRunLoopSourceRef src = CFRunLoopSourceCreate(NULL, 0, &ctx);
	if (src != NULL) {
		CFRunLoopAddSource(loop, src, kCFRunLoopCommonModes);
	}
	return src;
}

void releaseStopSource(CFRunLoopSourceRef src, CFRunLoopRef loop) {
	if (src != NULL && loop != NULL) {
		CFRunLoopRemoveSource(loop, src, kCFRunLoopCommonModes);
		CFRelease((CFTypeRef)src);
	}
}

void signalStopSource(CFRunLoopSourceRef src, CFRunLoopRef loop) {
	if (src != NULL) {
		CFRunLoopSourceSignal(src);
		if (loop != NULL) {
			CFRunLoopWakeUp(loop);
		}
	}
}
*/
import "C"

import (
	"runtime"
)

type RunLoop struct {
	loop       C.CFRunLoopRef
	stopSource C.CFRunLoopSourceRef
}

// CurrentRunLoop returns the current thread's CFRunLoop.
func CurrentRunLoop() *RunLoop {
	loop := C.currentRunLoop()
	C.retainRunLoop(loop)
	return &RunLoop{loop: loop, stopSource: 0}
}

// Run starts the CFRunLoop on the current thread.
func (r *RunLoop) Run() {
	if r == nil || r.loop == 0 {
		return
	}
	C.runLoopRun()
}

// StartRunLoop starts a CFRunLoop on a locked OS thread.
// The returned RunLoop can be stopped with Stop().
func StartRunLoop() *RunLoop {
	type rlData struct {
		loop C.CFRunLoopRef
		src  C.CFRunLoopSourceRef
	}
	ch := make(chan rlData, 1)
	go func() {
		runtime.LockOSThread()
		loop := C.currentRunLoop()
		C.retainRunLoop(loop)
		src := C.createStopSource(loop)
		ch <- rlData{loop: loop, src: src}
		C.runLoopRun()
		C.releaseStopSource(src, loop)
		C.releaseRunLoop(loop)
	}()
	d := <-ch
	return &RunLoop{loop: d.loop, stopSource: d.src}
}

// Stop stops the CFRunLoop.
func (r *RunLoop) Stop() {
	if r == nil || r.loop == 0 {
		return
	}
	if r.stopSource != 0 {
		C.signalStopSource(r.stopSource, r.loop)
		return
	}
	C.runLoopStop(r.loop)
}
