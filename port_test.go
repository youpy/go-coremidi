package coremidi

import "testing"
import "fmt"

func TestNewOutputPort(t *testing.T) {
	client, err := NewClient("test")

	if err != nil {
		t.Fatalf("failed to create a client")
	}

	_, err = NewOutputPort(client, "test")

	if err != nil {
		t.Fatalf("failed to create a port")
	}
}

func TestNewInputPort(t *testing.T) {
	client, _ := NewClient("test")

	port, err := NewInputPort(client, "test", func(value []byte) {
		fmt.Printf("value: %v\n", value)
	})

	if err != nil {
		t.Fatalf("failed to create a port")
	}

	sources, _ := AllSources()

	port.Connect(sources[0])

	packet := NewPacket([]byte{0x90, 0x30, 100})
	packet.Received(&sources[0])
}
