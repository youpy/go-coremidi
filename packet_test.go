package coremidi

import "testing"

func TestSend(t *testing.T) {
	devices, _ := AllDevices()
	device := devices[0]
	client, _ := NewClient("test client")
	port, _ := NewOutputPort(client, "test port")
	entities, _ := device.Entities()
	destinations, _ := entities[0].Destinations()
	destination := destinations[0]
	packet := NewPacket([]byte{0x90, 0x30, 100}, 0)

	err := packet.Send(&port, &destination)

	if err != nil {
		t.Fatalf("failed to send MIDI")
	}
}

func TestReceived(t *testing.T) {
	client, _ := NewClient("test client")
	source, _ := NewSource(client, "test source")
	packet := NewPacket([]byte{0x90, 0x30, 100}, 0)

	err := packet.Received(&source)

	if err != nil {
		t.Fatalf("failed to transmit MIDI")
	}
}
