package aggregate

import (
	"github.com/iwyg/go-eventsourced/identity"
	"github.com/iwyg/go-eventsourced/message"
	"sync"
)

type RootAggregate interface {
	ID() identity.ID
	Version() uint64
	Replay(<-chan message.Event) error
}

type Recorder interface {
	RecordThat(message.Event)
}

type EventFlusher interface {
	FlushEvents() <-chan message.Event
}

type Applicator interface {
	ApplyThat(message.Event)
}

type Root struct {
	mu             sync.RWMutex
	version        uint64
	Applicator     Applicator
	id             ID
	recordedEvents []message.Event
}

func (r *Root) Version() uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.version
}

func (r *Root) Replay(stream <- chan message.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for ev := range stream {
		r.Applicator.ApplyThat(ev)
		r.version = ev.AggregateVersion()
	}
	return nil
}


func (r *Root) ID() ID {
	return r.id
}

func (r *Root) FlushEvents() <-chan message.Event {
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

func (r *Root) RecordThat(ev message.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Applicator.ApplyThat(ev)
	r.version++
	event := ev.WithAggregateVersion(r.version)
	r.recordedEvents = append(r.recordedEvents, event)
}
