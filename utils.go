package emboxen

import (
	"io"
)

func OpenSourceCodeReader(ch *ControlChannel) io.Reader {
	rawInput := ch.inputChanel
	filteredInput := make(chan BuildRequestEvent, 16)
	ch.inputChanel = filteredInput

	reader, writer := io.Pipe()

	go func() {
		for {
			evt := <-rawInput
			if evt == nil {
				close(filteredInput)
				return
			}
			switch evt := evt.(type) {
			case *SourceCodeFragmentEvent:
				writer.Write(evt.Data)
			case *EndOfSourceCodeFragmentEvent:
				writer.Close()
			default:
				filteredInput <- evt
			}
		}
	}()

	return reader
}

/*
func openWriterGeneric(ch *ControlChannel) io.Writer {
	reader, writer := io.Pipe()

	go func() {
		buf := make([]byte, 16384)

		for {
			count, err := reader.Read(buf)
			// TODO: openWriterGeneric
		}
	}()

	return writer
}*/
