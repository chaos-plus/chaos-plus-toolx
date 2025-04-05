package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chaos-plus/chaos-plus-toolx/xsignal"
)

func main() {

	go func() {
		for i := 1; i < 10; i++ {
			fmt.Println("ticking..#####", i)
			time.Sleep(1 * time.Second)
		}
	}()
	xsignal.SetInteruptionSignal(func() {
		fmt.Println("Signal Interrupt 1")
		os.Exit(0)
	})
	xsignal.SetInteruptionSignal(func() {
		fmt.Println("Signal Interrupt 2")
	})
	for i := 1; i < 10; i++ {
		fmt.Println("ticking..@@@@@", i)
		time.Sleep(1 * time.Second)
	}
}
