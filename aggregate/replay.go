package aggregate

import "github.com/iwyg/go-eventsourced/message"

type RePlayer interface {
	// Replay replays inbound a stream events onto a target
	Replay(<-chan message.Event) error
}
