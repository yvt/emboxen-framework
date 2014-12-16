package emboxen

import (
	"encoding/gob"
	"os"
)

type controlChannel struct {
	outputChannel chan<- BuildEvent
	inputChanel   <-chan BuildRequestEvent
}

var _ctrlChannel controlChannel
var ctrlChannelInitialized bool

func ctrlChannel() *controlChannel {
	if ctrlChannelInitialized {
		return &_ctrlChannel
	}

	outputChannel := make(chan BuildEvent, 16)
	_ctrlChannel.outputChannel = outputChannel

	inputChannel := make(chan BuildRequestEvent, 16)
	_ctrlChannel.inputChanel = inputChannel

	go func() {
		g := gob.NewEncoder(os.Stdout)
		for {
			outputEvent := <-outputChannel
			if outputEvent == nil {
				return
			}

			err := g.Encode(outputEvent)

			if err != nil {
				panic(err)
			}
		}
	}()

	go func() {
		g := gob.NewDecoder(os.Stdin)
		for {
			var inputEvent BuildRequestEvent
			err := g.Decode(inputEvent)

			if err != nil {
				panic(err)
			}

			inputChannel <- inputEvent

			if _, ok := inputEvent.(EndOfSourceCodeFragmentEvent); ok {
				close(inputChannel)
				return
			}
		}
	}()

	ctrlChannelInitialized = true
	return &_ctrlChannel
}
