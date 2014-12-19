package emboxen

import "encoding/gob"

// This packages defines events sent from building
// environment to the Emboxen engine.

type BuildEvent interface{}

// Sent when buliding program has generated a message useful to
// the user.
type BuildOutputEvent struct {
	Output []byte
}

// Sent when building process has failed due to compilation error.
// This is one of the last event sent during a single connection.
type BuildFailedEvent struct{}

// Sent when building process has failed, but it was not because of the compiled program, but
// caused by the problem of the request.
// This is one of the last event sent during a single connection.
type BadRequestEvent struct {
	Error string
}

// Sent when building process has failed, but it was not because of the compiled program, but
// caused by an unexpected error.
// This is one of the last event sent during a single connection.
type SystemErrorEvent struct {
	Error string
}

// Sent when building process has suceeded, and
// the built program is about to be started in the building
// environment.
//
// After StartingProgramEvent is sent, neither of
// StartingProgramEvent and SendingProgramEvent must not be sent.
type StartingProgramEvent struct{}

// Sent when building process has suceeded, and
// the built program is about to be sent.
//
// After SendingProgramEvent is sent, neither of
// StartingProgramEvent and SendingProgramEvent must not be sent.
type SendingProgramEvent struct {
	Type ClientSideProgramType
}

// Sent after StartingProgramEvent was sent.
type ProgramOutputEvent struct {
	Output []byte
}

// Sent after SendingProgramEvent was sent.
type ProgramFragmentEvent struct {
	Data []byte
}

// Sent when executing the built program has completed.
// This is one of the last event sent during a single connection.
type EndOfProgramOutputEvent struct{}

// Sent after all parts of the built program has been sent.s
// This is one of the last event sent during a single connection.
type EndOfProgramFragmentEvent struct{}

type ClientSideProgramType uint32

const (
	_ ClientSideProgramType = iota
	ElfMipsLittleEndian
)

// Events below are sent by frontend.
type BuildRequestEvent interface{}

type SourceCodeFragmentEvent struct {
	Data []byte
}
type EndOfSourceCodeFragmentEvent struct{}

func init() {
	gob.Register(BuildOutputEvent{})
	gob.Register(BuildFailedEvent{})
	gob.Register(BadRequestEvent{})
	gob.Register(SystemErrorEvent{})
	gob.Register(StartingProgramEvent{})
	gob.Register(SendingProgramEvent{})
	gob.Register(ProgramOutputEvent{})
	gob.Register(ProgramFragmentEvent{})
	gob.Register(EndOfProgramOutputEvent{})
	gob.Register(EndOfProgramFragmentEvent{})
	gob.Register(SourceCodeFragmentEvent{})
	gob.Register(EndOfSourceCodeFragmentEvent{})
}
