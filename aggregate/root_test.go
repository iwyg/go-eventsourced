package aggregate

import (
	"github.com/iwyg/go-eventsourced/message"
	"sync"
	"testing"
)

type Account struct {
	Root
	sync.RWMutex
	balance int64
}

type AccountOpened struct {}
type MoneyAdded struct {Amount uint64}
type MoneyWithdrawn struct {Amount uint64}

func newModel() *Account {
	m  := newEmpty()
	id := NewID()
	m.Root.Applicator = m
	m.Record(message.NewEvent(message.ID(id), &AccountOpened{}))
	return m
}

func newEmpty() *Account {
	m := &Account{}
	m.Root.Applicator = m
	return m
}

func (m *Account) Balance() int64 {
	m.RLock()
	defer m.RUnlock()
	return m.balance
}

func (m *Account) Apply(event message.Event) {
	switch event.Type().(type) {
		case *AccountOpened:
			m.Root.WithID(ID(event.AggregateID()))
			break
		case *MoneyAdded:
			ev, _ := event.Type().(*MoneyAdded)
			m.balance = m.balance + int64(ev.Amount)
			break
		case *MoneyWithdrawn:
			ev, _ := event.Type().(*MoneyWithdrawn)
			m.balance = m.balance - int64(ev.Amount)
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
	m := newEmpty()
	events := []message.Event{
		message.NewEvent(message.ID(id), &AccountOpened{}).WithAggregateVersion(1),
		message.NewEvent(message.ID(id), &MoneyAdded{Amount: 100}).WithAggregateVersion(2),
		message.NewEvent(message.ID(id), &MoneyWithdrawn{Amount: 34}).WithAggregateVersion(3),
	}

	for _, ev := range events {
		m.Record(ev)
	}

	if len(m.recordedEvents) != 3 {
		t.Errorf("want 3, got %d", len(m.recordedEvents))
	}

	_ = m.Flush()

	if len(m.recordedEvents) != 0 {
		t.Errorf("want 0, got %d", len(m.recordedEvents))
	}
}


func TestRoot_replay_panic_nil_applicator(t *testing.T) {
	defer func() {
		r := recover()

		if r == nil {
			t.Error("should have caused panic")
		}

		if r != "nil Root.Applicator" {
			t.Errorf("expected %q, got %q", "nil Root.Applicator", r)
		}
	}()

	id := NewID()
	m  := &Account{}
	events := []message.Event{
		message.NewEvent(message.ID(id), &AccountOpened{}).WithAggregateVersion(1),
	}

	stream := make(chan message.Event, len(events))
	go func () {
		defer close(stream)
		for _, ev := range events {
			stream <- ev
		}
	}()

	if err := m.Replay(stream); err != nil {
		t.Error(err)
	}
}


func TestRoot_Record_panic_nil_applicator(t *testing.T) {
	defer func() {
		r := recover()

		if r == nil {
			t.Error("should have caused panic")
		}

		if r != "nil Root.Applicator" {
			t.Errorf("expected %q, got %q", "nil Root.Applicator", r)
		}
	}()

	m := &Account{}
	m.Record(message.NewEvent(message.ID(m.ID()), &AccountOpened{}).WithAggregateVersion(1))
}

func TestRoot_replay_entity(t *testing.T) {
	id := NewID()
	m := newEmpty()
	events := []message.Event{
		message.NewEvent(message.ID(id), &AccountOpened{}).WithAggregateVersion(1),
		message.NewEvent(message.ID(id), &MoneyAdded{180}).WithAggregateVersion(2),
		message.NewEvent(message.ID(id), &MoneyWithdrawn{32}).WithAggregateVersion(3),
		message.NewEvent(message.ID(id), &MoneyAdded{223}).WithAggregateVersion(4),
		message.NewEvent(message.ID(id), &MoneyWithdrawn{23}).WithAggregateVersion(5),
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

	if m.Version() != 5 {
		t.Errorf("want model version 3, got %d", m.Version())
	}

	var expectedBalance int64 = 348
	if m.Balance() != expectedBalance {
		t.Errorf("want %d, got %d", expectedBalance, m.Balance())
	}
}
