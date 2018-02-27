package coremidi

import (
	"bytes"
	"testing"
	"time"
)

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
	ch := make(chan Packet)

	port, err := NewInputPort(client, "test", func(source Source, packet Packet) {
		ch <- packet
	})

	if err != nil {
		t.Fatalf("failed to create a port")
	}

	sources, _ := AllSources()

	connection, _ := port.Connect(sources[0])

	packet := NewPacket([]byte{0x90, 0x30, 100}, 2345)
	packet.Received(&sources[0])

	select {
	case packet := <-ch:
		if bytes.Compare(packet.Data, []byte{0x90, 0x30, 100}) != 0 {
			t.Fatalf("invalid value: %v", packet.Data)
		}

		if packet.TimeStamp != 2345 {
			t.Fatalf("invalid timestamp: %v", packet.TimeStamp)
		}

		connection.Disconnect()
	case <-time.After(1 * time.Second):
		t.Fatal("timed out")
	}
}
