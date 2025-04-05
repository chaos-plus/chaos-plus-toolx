package xsignal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func SetInteruptionSignal(hook func()) chan<- os.Signal {
	if hook == nil {
		return nil
	}
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r\nCtrl+C pressed in Terminal") 
		hook()
	}()
	return c
}
