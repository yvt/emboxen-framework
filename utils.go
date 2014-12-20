package emboxen

import (
	"io"
)

func markEndOfOutput(ch *controlChannel) {
	go func() {
		close(ch.outputChannel)
	}()
	<-ch.waitForFlush
}

func openSourceCodeReader(ch *controlChannel) io.ReadCloser {
	rawInput := ch.inputChanel
	filteredInput := make(chan BuildRequestEvent, 16)
	ch.inputChanel = filteredInput

	reader, writer := io.Pipe()
	closed := false

	go func() {
		for {
			evt := <-rawInput
			if evt == nil {
				close(filteredInput)
				if !closed {
					panic("Unexpected closure.")
				}
				return
			}
			switch evt := evt.(type) {
			case SourceCodeFragmentEvent:
				_, err := writer.Write(evt.Data)
				if err != nil {
					// Unexpected closure.
					panic("Unexpected closure.")
				}
			case EndOfSourceCodeFragmentEvent:
				writer.Close()
				closed = true
			default:
				filteredInput <- evt
			}
		}
	}()

	return reader
}

func openWriterGeneric(ch *controlChannel, sendChunk func([]byte), endOutput func()) io.WriteCloser {
	reader, writer := io.Pipe()

	go func() {
		buf := make([]byte, 16384)

		for {
			count, err := reader.Read(buf)
			if err == io.EOF {
				reader.Close()
				endOutput()
				break
			}
			if count > 0 {
				sendChunk(buf[0:count])
			}
		}
	}()

	return writer
}

func openProgramWriter(ch *controlChannel) io.WriteCloser {
	return openWriterGeneric(ch, func(p []byte) {
		ch.outputChannel <- ProgramFragmentEvent{p}
	}, func() {
		ch.outputChannel <- EndOfProgramFragmentEvent{}
	})
}

// This one does not send "end of ~~~" event.
func openBuildOutputWriter(ch *controlChannel) io.WriteCloser {
	return openWriterGeneric(ch, func(p []byte) {
		ch.outputChannel <- BuildOutputEvent{p}
	}, func() {})
}

func openApplicationOutputWriter(ch *controlChannel) io.WriteCloser {
	return openWriterGeneric(ch, func(p []byte) {
		ch.outputChannel <- ProgramOutputEvent{p}
	}, func() {
		ch.outputChannel <- EndOfProgramOutputEvent{}
	})
}
