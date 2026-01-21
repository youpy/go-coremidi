package coremidi

import "testing"

func TestRunLoopStartStop(t *testing.T) {
	loop := StartRunLoop()
	if loop == nil || loop.loop == 0 {
		t.Fatalf("expected valid run loop")
	}
	loop.Stop()
}

func TestCurrentRunLoop(t *testing.T) {
	loop := CurrentRunLoop()
	if loop == nil || loop.loop == 0 {
		t.Fatalf("expected current run loop")
	}
}

func TestRunLoopStopNil(t *testing.T) {
	var loop *RunLoop
	loop.Stop()
}
