package coremidi

import "testing"

func TestNotificationMessageIDString(t *testing.T) {
	cases := []struct {
		id       MIDINotificationMessageID
		expected string
	}{
		{MIDIMsgSetupChanged, "SetupChanged"},
		{MIDIMsgObjectAdded, "ObjectAdded"},
		{MIDIMsgObjectRemoved, "ObjectRemoved"},
		{MIDIMsgPropertyChanged, "PropertyChanged"},
		{MIDIMsgThruConnectionsChanged, "ThruConnectionsChanged"},
		{MIDIMsgSerialPortOwnerChanged, "SerialPortOwnerChanged"},
		{MIDIMsgIOError, "IOError"},
		{MIDINotificationMessageID(999), "Unknown"},
	}

	for _, testCase := range cases {
		if got := testCase.id.String(); got != testCase.expected {
			t.Fatalf("unexpected string for %d: %s", testCase.id, got)
		}
	}
}
