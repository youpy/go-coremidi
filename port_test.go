package coremidi

import "testing"
import "bytes"
import "time"

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
	ch := make(chan []byte)

	port, err := NewInputPort(client, "test", func(source Source, value []byte) {
		ch <- value
	})

	if err != nil {
		t.Fatalf("failed to create a port")
	}

	sources, _ := AllSources()

	connection, _ := port.Connect(sources[0])

	packet := NewPacket([]byte{0x90, 0x30, 100})
	packet.Received(&sources[0])

	select {
	case value := <-ch:
		if bytes.Compare(value, []byte{0x90, 0x30, 100}) != 0 {
			t.Fatalf("invalid value: %v", value)
		}

		connection.Disconnect()
	case <-time.After(1 * time.Second):
		t.Fatal("timed out")
	}
}
