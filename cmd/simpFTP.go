package main

import (
	"log"
	"simpFTP/internal"
)

func main() {
	var err error

	defer func() {
		if err != nil {
			log.Println(err)

		}
	}()

	go func() {
		err = internal.SimpServer.ListenAndServe()
	}()
	internal.MainLoop(internal.SimpServer)
}
