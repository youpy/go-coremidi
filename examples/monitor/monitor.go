package main

import (
	"fmt"

	coremidi "github.com/youpy/go-coremidi"
)

func main() {
	client, err := coremidi.NewClient("a client")
	if err != nil {
		fmt.Println(err)
		return
	}

	port, err := coremidi.NewInputPort(
		client,
		"test",
		func(source coremidi.Source, packet coremidi.Packet) {
			fmt.Printf(
				"device: %v, manufacturer: %v, source: %v, data: %v\n",
				source.Entity().Device().Name(),
				source.Manufacturer(),
				source.Name(),
				packet.Data,
			)
			return
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	sources, err := coremidi.AllSources()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, source := range sources {
		func(source coremidi.Source) {
			port.Connect(source)
		}(source)
	}

	ch := make(chan int)
	<-ch
}
