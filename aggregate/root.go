package aggregate

import (
	"github.com/iwyg/go-eventsourced/identity"
	"github.com/iwyg/go-eventsourced/message"
	"sync"
)

type Aggregate interface {
	// ID return the aggregate identity
	ID() identity.ID
	// Version return the aggregate version
	Version() uint64
}

type Recorder interface {
	// Records records an inbound event
	Record(message.Event)
}

type Flusher interface {
	// Flush flushes stored events and returns them as an outbound stream
	Flush() <-chan message.Event
}

type Replayer interface {
	// Replay replays inbound a stream events onto a target
	Replay(<-chan message.Event) error
}

type Applicator interface {
	// Apply applies an inbound event onto a target
	Apply(message.Event)
}

type Root struct {
	mu             sync.RWMutex
	version        uint64
	Applicator     Applicator
	id             ID
	recordedEvents []message.Event
}

// ID return the aggregate identity
func (r *Root) ID() ID {
	return r.id
}

func (r *Root) WithID(id ID) {
	r.id = id
}

// Version return the aggregate version
func (r *Root) Version() uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.version
}

// Replay replays inbound a stream events onto a target
func (r *Root) Replay(stream <- chan message.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Applicator == nil {
		panic("nil Root.Applicator")
		return nil
	}
	for ev := range stream {
		r.Applicator.Apply(ev)
		r.version = ev.AggregateVersion()
	}
	return nil
}


// Flush flushes stored events and returns them as an outbound stream
func (r *Root) Flush() <-chan message.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	flushed := make(chan message.Event, len(r.recordedEvents))
	go func (events []message.Event) {
		defer close(flushed)
		for _, e := range events {
			flushed <-e
		}
	}(r.recordedEvents)
	r.recordedEvents = r.recordedEvents[0:0]
	return flushed
}

// Records records an inbound event
func (r *Root) Record(ev message.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Applicator == nil {
		panic("nil Root.Applicator")
		return
	}
	r.Applicator.Apply(ev)
	r.version++
	event := ev.WithAggregateVersion(r.version)
	r.recordedEvents = append(r.recordedEvents, event)
}
