#include <CoreMIDI/CoreMIDI.h>
#include <stdint.h>

extern void midiNotifyCallback(uintptr_t handle, MIDINotification *message);

void midiNotifyProc(const MIDINotification *message, void *refCon) {
  midiNotifyCallback((uintptr_t)refCon, (MIDINotification *)message);
}
