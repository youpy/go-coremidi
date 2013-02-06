package coremidi

import "testing"
import "io"

func TestWrite(t *testing.T) {
	var writer io.Writer

	client, _ := NewClient("test client")
	port, _ := NewOutputPort(client, "test port")
	destinations, _ := AllDestinations()

	writer = &PacketWriter{&port, &destinations[0]}
	n, err := writer.Write([]byte{0x90, 0x30, 100})

	if err != nil {
		t.Fatalf("failed to write")
	}

	if n != 3 {
		t.Fatalf("invalid result")
	}
}
