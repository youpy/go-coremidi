#include <CoreMIDI/CoreMIDI.h>
#include <stdint.h>

extern void midiNotifyCallback(uintptr_t handle, MIDINotification *message);

void midiNotifyProc(const MIDINotification *message, void *refCon) {
  // refCon points to a C-allocated uintptr_t holding a Go cgo.Handle value.
  // See the Client.refCon field in client.go for why we box the handle rather
  // than casting it directly into the void* refCon slot.
  uintptr_t handle = refCon != NULL ? *(uintptr_t *)refCon : 0;
  midiNotifyCallback(handle, (MIDINotification *)message);
}
