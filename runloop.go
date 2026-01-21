package coremidi

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>

static CFRunLoopRef currentRunLoop() {
  return CFRunLoopGetCurrent();
}

static void retainRunLoop(CFRunLoopRef loop) {
  if (loop != NULL) {
    CFRetain(loop);
  }
}

static void releaseRunLoop(CFRunLoopRef loop) {
  if (loop != NULL) {
    CFRelease(loop);
  }
}

static void runLoopRun() {
  CFRunLoopRun();
}

static void runLoopStop(CFRunLoopRef loop) {
  if (loop != NULL) {
    CFRunLoopStop(loop);
    CFRunLoopWakeUp(loop);
  }
}
*/
import "C"

import "runtime"

type RunLoop struct {
	loop C.CFRunLoopRef
}

// CurrentRunLoop returns the current thread's CFRunLoop.
func CurrentRunLoop() *RunLoop {
	loop := C.currentRunLoop()
	C.retainRunLoop(loop)
	return &RunLoop{loop: loop}
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
	ch := make(chan C.CFRunLoopRef, 1)
	go func() {
		runtime.LockOSThread()
		loop := C.currentRunLoop()
		C.retainRunLoop(loop)
		ch <- loop
		C.runLoopRun()
		C.releaseRunLoop(loop)
	}()
	loop := <-ch
	return &RunLoop{loop: loop}
}

// Stop stops the CFRunLoop.
func (r *RunLoop) Stop() {
	if r == nil || r.loop == 0 {
		return
	}
	C.runLoopStop(r.loop)
}
