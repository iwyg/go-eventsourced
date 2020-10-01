package aggregate

import (
	"github.com/iwyg/go-eventsourced/message"
	"reflect"
	"sync"
	"testing"
)

type Model struct {
	Root
	sync.RWMutex
	a, b bool
}

type ModelCreated struct {}
type ModelEventA struct {}
type ModelEventB struct {}

func newModel() *Model {
	m := newEmpty()
	m.id = NewID()
	m.Root.Applicator = m
	m.RecordThat(message.NewEvent(message.ID(m.id), &ModelCreated{}))
	return m
}

func newEmpty() *Model {
	m := &Model{}
	m.Root.Applicator = m
	return m
}

func (m *Model) ReceivedA() bool {
	m.RLock()
	defer m.RUnlock()
	return m.a
}

func (m *Model) ReceivedB() bool {
	m.RLock()
	defer m.RUnlock()
	return m.b
}
func (m *Model) ApplyThat(event message.Event) {
	switch event.Type().(type) {
		case *ModelCreated:
			m.id = ID(event.AggregateID())
			break
		case *ModelEventA:
			m.a = true
			break
		case *ModelEventB:
			m.b = true
			break
	}
}

func TestRoot_in_entity_struct(t *testing.T) {
	m := newModel()

	if m.Version() != 1 {
		t.Errorf("want version 1 got %d", m.Version())
	}
}

func TestRoot_flush_events(t *testing.T) {
	id := NewID()
	m := newModel()
	events := []message.Event{
		message.NewEvent(message.ID(id), &ModelCreated{}).WithAggregateVersion(1),
		message.NewEvent(message.ID(id), &ModelEventA{}).WithAggregateVersion(2),
		message.NewEvent(message.ID(id), &ModelEventB{}).WithAggregateVersion(3),
	}

	for _, ev := range events {
		m.RecordThat(ev)
	}

	for range m.FlushEvents() {}

	if len(m.recordedEvents) != 0 {
		t.Errorf("want 0, got %d", len(m.recordedEvents))
	}
}

func TestRoot_replay_entity(t *testing.T) {
	id := NewID()
	m := newModel()
	events := []message.Event{
		message.NewEvent(message.ID(id), &ModelCreated{}).WithAggregateVersion(1),
		message.NewEvent(message.ID(id), &ModelEventA{}).WithAggregateVersion(2),
		message.NewEvent(message.ID(id), &ModelEventB{}).WithAggregateVersion(3),
	}

	stream := make(chan message.Event, len(events))
	go func () {
		defer close(stream)
		for _, ev := range events {
			stream <- ev
		}
	}()

	if err := m.Replay(stream); err != nil {
		t.Fatal(err)
	}

	if !m.ID().Equals(id) {
		t.Errorf("want model id %s, got %s", id, m.ID())
	}

	if m.Version() != 3 {
		t.Errorf("want model version 3, got %d", m.Version())
	}

	if !m.ReceivedA() {
		t.Errorf("model should've received %s", reflect.TypeOf(&ModelEventA{}))
	}

	if !m.ReceivedB() {
		t.Errorf("model should've received %s", reflect.TypeOf(&ModelEventB{}))
	}
}
