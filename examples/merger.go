package main

import (
	coremidi "github.com/youpy/go-coremidi"
)

func main() {
	client, err := coremidi.NewClient("a client")
	if err != nil {
		panic(err)
	}

	sources, err := coremidi.AllSources()
	if err != nil {
		panic(err)
	}

	mixerSource, err := coremidi.NewSource(client, "FM source")
	if err != nil {
		panic(err)
	}

	mixerDestination, err := coremidi.NewDestination(client, "FM destination", func(packet coremidi.Packet) {
		err := packet.Received(&mixerSource)
		if err != nil {
			panic(err)
		}
		return
	})
	if err != nil {
		panic(err)
	}

	outPort, err := coremidi.NewOutputPort(client, "FM in")
	if err != nil {
		panic(err)
	}
	port, err := coremidi.NewInputPort(client, "test", func(source coremidi.Source, packet coremidi.Packet) {
		packet.Send(&outPort, &mixerDestination)
		return
	})

	if err != nil {
		panic(err)
	}

	for _, source := range sources {
		if source.Name() != mixerSource.Name() {
			func(source coremidi.Source) {
				port.Connect(source)
			}(source)
		}
	}

	ch := make(chan int)
	<-ch
}
