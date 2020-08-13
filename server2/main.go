package main

import (
	"github.com/hashlash/descarca/pipe"
	"log"
)

func main() {
	p, err := pipe.SetupDownloadJobPipe()
	if err != nil {
		log.Fatal("Cannot create amqp pipe: ", err)
	}

	consumer, err := p.Receive()
	if err != nil {
		log.Fatal("Cannot retrieve consumer channel: ", err)
	}

	log.Println("Waiting for download job")

	for data := range consumer {
		n, err := download(data["key"], data["url"], data["uuid"])
		if err != nil {
			log.Println("Error while downloading: ", err)
			continue
		}
		log.Printf("Successfully download %v bytes to disk\n", n)
	}
}
