package main

import (
	"flag"
	"fmt"
)

var mode int

const (
	server = iota + 1
	client
)

func main() {
	flag.IntVar(&mode, "mode", 1, "mode = 1 for server, 2 for client")
	flag.Parse()

	switch mode {
	case server:
		fmt.Println("server")
	case client:
		fmt.Println("client")
	default:
		return
	}
}
