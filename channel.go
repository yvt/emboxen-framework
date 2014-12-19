package emboxen

import (
	"bytes"
	"encoding/gob"
	"io"
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
		buf := &bytes.Buffer{}
		cout := os.Stdout
		lenbuf := [4]byte{}
		for {
			outputEvent := <-outputChannel
			if outputEvent == nil {
				return
			}

			buf.Reset()
			buf.Write(lenbuf[0:]) // placeholder for chunk length
			g := gob.NewEncoder(buf)
			err := g.Encode(&outputEvent)
			if err != nil {
				panic(err)
			}

			p := buf.Bytes()
			ln := len(p) - 4
			p[0] = byte(ln)
			p[1] = byte(ln >> 8)
			p[2] = byte(ln >> 16)
			p[3] = byte(ln >> 24)
			_, err = cout.Write(p)
			if err != nil {
				panic(err)
			}
		}
	}()

	go func() {
		lenbuf := [4]byte{}
		cin := os.Stdin
		for {
			_, err := io.ReadFull(cin, lenbuf[0:])
			if err != nil {
				panic(err)
			}

			ln := int(lenbuf[0]) | int(lenbuf[1])<<8 | int(lenbuf[2])<<16 | int(lenbuf[3])<<24

			var inputEvent BuildRequestEvent
			g := gob.NewDecoder(io.LimitReader(cin, int64(ln)))
			err = g.Decode(&inputEvent)

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
