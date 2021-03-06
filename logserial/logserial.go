package main

import (
	"bufio"
	"log"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: DEVICE, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(s)
	for {
		reply, err := reader.ReadBytes('\x0a')
		if err != nil {
			panic(err)
		}
		log.Print(string(reply))
	}

	s.Close()
}
