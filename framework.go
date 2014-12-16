package emboxen

import (
	"fmt"
	"io"
	"log"
	"os"
)

type InternalErrorReporter interface {
	ReportInternalError(e interface{})
}

type SourceCodeReader interface {
	InternalErrorReporter
	io.Reader
	CloseAndOpenBuildOuputWriter() BuildOutputWriter
}

type BuildOutputWriter interface {
	InternalErrorReporter
	io.Writer

	CloseWithFailure()
	CloseAndOpenProgramWriter(programType ClientSideProgramType) ProgramWriter
	CloseAndOpenApplicationOutputWriter() ApplicationOutputWriter
}

type ProgramWriter interface {
	InternalErrorReporter
	io.WriteCloser
}

type ApplicationOutputWriter interface {
	InternalErrorReporter
	io.WriteCloser
}

type HandlerFunc func(SourceCodeReader)

func Accept(handler HandlerFunc) {
	log.SetOutput(os.Stderr)
	handler(&sourceCodeReader{
		ch:     ctrlChannel(),
		reader: openSourceCodeReader(ctrlChannel()),
	})
}

type sourceCodeReader struct {
	reader io.ReadCloser
	ch     *controlChannel
}

func (self *sourceCodeReader) Read(p []byte) (n int, err error) {
	if self.reader == nil {
		panic("Attempted to read SourceCodeReader after closing it.")
	}
	n, err = self.reader.Read(p)
	return
}

func (self *sourceCodeReader) ReportInternalError(e interface{}) {
	if self.reader == nil {
		panic("Attempted to close SourceCodeReader twice.")
	}
	self.reader.Close()
	self.reader = nil
	self.ch.outputChannel <- SystemErrorEvent{
		Error: fmt.Sprintf("%v", e),
	}
}

func (self *sourceCodeReader) CloseAndOpenBuildOuputWriter() BuildOutputWriter {
	if self.reader == nil {
		panic("Attempted to close SourceCodeReader twice.")
	}
	self.reader.Close()
	self.reader = nil
	return &buildOutputWriter{
		writer: openBuildOutputWriter(self.ch),
		ch:     self.ch,
	}
}

type buildOutputWriter struct {
	writer io.WriteCloser
	ch     *controlChannel
}

func (self *buildOutputWriter) Write(p []byte) (n int, err error) {
	if self.writer == nil {
		panic("Attempted to write to BuildOutputWriter after closing it.")
	}
	n, err = self.writer.Write(p)
	return
}
func (self *buildOutputWriter) innerClose() {
	if self.writer == nil {
		panic("Attempted to close BuildOutputWriter twice.")
	}
	self.writer.Close()
	self.writer = nil
}
func (self *buildOutputWriter) ReportInternalError(e interface{}) {
	if self.writer == nil {
		panic("Attempted to close BuildOutputWriter twice.")
	}
	self.ch.outputChannel <- SystemErrorEvent{
		Error: fmt.Sprintf("%v", e),
	}
	self.writer.Close()
	self.writer = nil
}
func (self *buildOutputWriter) CloseWithFailure() {
	self.innerClose()
	self.ch.outputChannel <- BuildFailedEvent{}
}
func (self *buildOutputWriter) CloseAndOpenProgramWriter(programType ClientSideProgramType) ProgramWriter {
	self.innerClose()
	self.ch.outputChannel <- SendingProgramEvent{
		Type: programType,
	}
	return &genericWriteCloser{
		writer: openProgramWriter(self.ch),
		ch:     self.ch,
	}
}
func (self *buildOutputWriter) CloseAndOpenApplicationOutputWriter() ApplicationOutputWriter {
	self.innerClose()
	self.ch.outputChannel <- StartingProgramEvent{}
	return &genericWriteCloser{
		writer: openApplicationOutputWriter(self.ch),
		ch:     self.ch,
	}
}

type genericWriteCloser struct {
	writer io.WriteCloser
	ch     *controlChannel
}

func (self *genericWriteCloser) Write(p []byte) (n int, err error) {
	if self.writer == nil {
		panic("Attempted to write after closing it.")
	}
	n, err = self.writer.Write(p)
	return
}
func (self *genericWriteCloser) Close() error {
	if self.writer == nil {
		panic("Attempted to close twice.")
	}
	r := self.writer.Close()
	self.writer = nil
	return r
}
func (self *genericWriteCloser) ReportInternalError(e interface{}) {
	if self.writer == nil {
		panic("Attempted to close twice.")
	}
	self.ch.outputChannel <- SystemErrorEvent{
		Error: fmt.Sprintf("%v", e),
	}
	self.writer.Close()
	self.writer = nil
}
