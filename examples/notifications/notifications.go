package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	coremidi "github.com/youpy/go-coremidi"
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
