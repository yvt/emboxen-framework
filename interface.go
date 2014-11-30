package emboxen

// This packages defines events sent from building
// environment to the Emboxen engine.

type BuildEvent interface{}

// Sent when buliding program has generated a message useful to
// the user.
type BuildOutputEvent struct {
	Segments []TextSegment
}

// Sent when building process has failed due to compilation error.
// This is one of the last event sent during a single connection.
type BuildFailedEvent struct{}

// Sent when building process has failed, but it was not because of the compiled program.
// This is one of the last event sent during a single connection.
type BadRequestEvent struct {
	Message string
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
	Segments []TextSegment
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

type TextSegment struct {
	Text  string
	Red   uint8
	Green uint8
	Blue  uint8
	Flags TextSegmentFlags
}

type TextSegmentFlags uint8

const (
	UseColor TextSegmentFlags = 1 << iota
	Bold                      = 1 << iota
)
