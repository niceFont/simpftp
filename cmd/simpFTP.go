package main

import (
	"log"
	"simpFTP/internal"
)

func main() {
	var err error

	defer func() {
		log.Fatal(err)
	}()

	go func() {
		err = internal.SimpServer.ListenAndServe()
	}()
	internal.MainLoop(internal.SimpServer)
}
