package aggregate

import "github.com/iwyg/go-eventsourced/message"

type Recorder interface {
	// Records records an inbound event
	Record(message.Event)
}

type Flusher interface {
	// Flush flushes stored events and returns them as an outbound stream
	Flush() <-chan message.Event
}

type Applicator interface {
	// Apply applies an inbound event onto a target
	Apply(message.Event)
}

