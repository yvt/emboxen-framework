package emboxen

import (
	"log"
	"os"
)

type HandlerFunc func(*ControlChannel)

func Accept(handler HandlerFunc) {
	log.SetOutput(os.Stderr)
	handler(ctrlChannel())
}
