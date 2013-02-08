# go-coremidi

A Go library using MIDI on Mac

## Installation

```
go get github.com/youpy/go-coremidi
```

## Synopsis

```go
package main

import (
	"fmt"
	"github.com/youpy/go-coremidi"
	"time"
)

func main() {
	client, err := coremidi.NewClient("a client")

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
			port, err := coremidi.NewInputPort(client, "test", func(value []byte) {
				fmt.Printf("source: %v manufacturer: %v value: %v\n", source.Name(), source.Manufacturer(), value)
				return
			})

			if err != nil {
				fmt.Println(err)
				return
			}

			port.Connect(source)
		}(source)
	}

	time.Sleep(10 * time.Second)
}
```
