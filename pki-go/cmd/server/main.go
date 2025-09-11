package main

import (
	"log"

	"github.com/sirajudheenam/pki/pki-go/internal/server"
)

func main() {
	srv, err := server.NewServer(":8443", "./certs/server")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.Start())
}
