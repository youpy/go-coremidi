# go-coremidi

A Go library to use MIDI on Mac

## Installation

```
go get github.com/youpy/go-coremidi
```

## Synopsis

### Monitor MIDI Messages

```go
package main

import (
	"fmt"

	"github.com/youpy/go-coremidi"
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
```

### Handle Device Notifications

CoreMIDI delivers device notifications via a CFRunLoop on the same OS thread
that created the client. If the run loop is not running, device changes may
not be delivered and polling APIs can return stale results.

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/youpy/go-coremidi"
)

func main() {
	runtime.LockOSThread()

	loop := coremidi.CurrentRunLoop()
	client, err := coremidi.NewClientWithNotification("notify-client", func(notification coremidi.Notification) {
		fmt.Printf("notification: %s\n", notification.MessageID)
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		<-stop
		loop.Stop()
	}()

	loop.Run()
}
```

## Documents

* http://godoc.org/github.com/youpy/go-coremidi
