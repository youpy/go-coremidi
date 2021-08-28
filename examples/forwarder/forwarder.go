// Forward incoming MIDI messages from MIDI keyboard to another MIDI interface

package main

import (
	"flag"

	coremidi "github.com/youpy/go-coremidi"
)

func main() {
	var sourceName string
	var destinationName string

	flag.StringVar(&sourceName, "source", "KEYBOARD", "source name")
	flag.StringVar(&destinationName, "destination", "UM-ONE", "destination name")
	flag.Parse()

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

	var targetDestination coremidi.Destination
	for _, destination := range destinations {
		if destination.Name() == destinationName {
			targetDestination = destination
		}
	}

	outPort, err := coremidi.NewOutputPort(client, "an output")
	if err != nil {
		panic(err)
	}

	port, err := coremidi.NewInputPort(client, "an input", func(source coremidi.Source, packet coremidi.Packet) {
		packet.Send(&outPort, &targetDestination)
	})
	if err != nil {
		panic(err)
	}

	for _, source := range sources {
		if source.Name() == sourceName {
			port.Connect(source)
		}
	}

	ch := make(chan int)
	<-ch
}
