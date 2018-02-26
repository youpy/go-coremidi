package coremidi

import (
	"bytes"
	"testing"
	"time"
)

func TestNewDestination(t *testing.T) {
	client, err := NewClient("a client")
	if err != nil {
		panic(err)
	}

	sources, err := AllSources()
	if err != nil {
		panic(err)
	}

	mixerSource, err := NewSource(client, "FM source")
	if err != nil {
		panic(err)
	}

	ch := make(chan []byte)
	mixerDestination, err := NewDestination(client, "FM destination", func(value []byte) {
		err := NewPacket(value).Received(&mixerSource)
		if err != nil {
			panic(err)
		}

		ch <- value

		return
	})
	if err != nil {
		panic(err)
	}

	outPort, err := NewOutputPort(client, "FM in")
	if err != nil {
		panic(err)
	}
	port, err := NewInputPort(client, "test", func(source Source, value []byte) {
		NewPacket(value).Send(&outPort, &mixerDestination)
		return
	})

	if err != nil {
		panic(err)
	}

	for _, source := range sources {
		if source.Name() != mixerSource.Name() {
			func(source Source) {
				port.Connect(source)
			}(source)
		}
	}

	packet := NewPacket([]byte{0x90, 0x30, 100})
	packet.Received(&sources[0])

	select {
	case value := <-ch:
		if bytes.Compare(value, []byte{0x90, 0x30, 100}) != 0 {
			t.Fatalf("invalid value: %v", value)
		}

		mixerDestination.Dispose()
	case <-time.After(1 * time.Second):
		t.Fatal("timed out")
	}
}

func TestNumberOfDestinations(t *testing.T) {
	destinations, err := AllDestinations()
	numberOfDestinations := len(destinations)

	if err != nil {
		t.Fatalf("failed to get destinations")
	}

	if numberOfDestinations <= 0 {
		t.Fatalf("invalid number of destinations")
	}
}
