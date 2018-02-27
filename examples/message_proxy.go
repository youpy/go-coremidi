// Proxy incoming MIDI messages from MIDI keyboard to another MIDI interface

package main

import (
	coremidi "github.com/youpy/go-coremidi"
)

func main() {
	var targetDestination coremidi.Destination

	client, err := coremidi.NewClient("a client")
	if err != nil {
		panic(err)
	}

	sources, err := coremidi.AllSources()
	if err != nil {
		panic(err)
	}

	destinations, err := coremidi.AllDestinations()
	if err != nil {
		panic(err)
	}

	for _, destination := range destinations {
		if destination.Name() == "UM-ONE" {
			targetDestination = destination
		}
	}

	outPort, err := coremidi.NewOutputPort(client, "an output")
	if err != nil {
		panic(err)
	}

	port, err := coremidi.NewInputPort(client, "an input", func(source coremidi.Source, packet coremidi.Packet) {
		packet.Send(&outPort, &targetDestination)

		return
	})
	if err != nil {
		panic(err)
	}

	for _, source := range sources {
		if source.Name() == "KEYBOARD" {
			port.Connect(source)
		}
	}

	ch := make(chan int)
	<-ch
}
