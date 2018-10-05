package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: DEVICE, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}

	cc := make(chan os.Signal, 2)
	signal.Notify(cc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cc
		s.Close()
		fmt.Printf("\n^C detected, exit.\n")
		os.Exit(0)
	}()

	b := make([]byte, 1)
	for {
		_, err := s.Read(b)
		if err != nil {
			panic(err)
		}
		fmt.Print(string(b))
	}
}
