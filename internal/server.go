package internal

import (
	"net"
	"net/http"

	"log"
)

//SimpServer uses the http package server to create a server instance
var SimpServer *http.Server

func init() {

	myLogger := &log.Logger{}

	SimpServer = &http.Server{
		Addr:              "localhost:3000",
		Handler:           http.FileServer(http.Dir("./www")),
		TLSConfig:         nil,
		ReadTimeout:       30,
		ReadHeaderTimeout: 0,
		WriteTimeout:      30,
		IdleTimeout:       30,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState: func(connection net.Conn, state http.ConnState) {
			switch state {

			case http.StateNew:
				log.Println("Request Incoming")
			case http.StateActive:
				log.Println("Handling Request")
			case http.StateIdle:
				log.Println("Keep-Alive")
			case http.StateClosed:
				log.Println("Connection closed.")
			}

		},
		ErrorLog: myLogger,
	}

}
