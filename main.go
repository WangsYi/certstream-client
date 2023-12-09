package main

import (
	"github.com/WangsYi/certstream-go"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

func main() {
	// The false flag specifies that we want heartbeat messages.
	stream, errStream := certstream.CertStreamEventStream(false, "wss://127.0.0.1:5000/full_stream")
	for {
		select {
		case jq := <-stream:
			messageType, err := jq.String("message_type")

			if err != nil {
				log.Fatal("Error decoding jq string")
			}

			log.Info("Message type -> ", messageType)
			log.Info("recv: ", jq)

		case err := <-errStream:
			log.Error(err)
		}
	}
}
