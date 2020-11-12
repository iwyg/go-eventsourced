package message

import (
	"encoding/json"
	"time"
)

const (
	metaKeyAggVers MetaKey = "@@aggregate_id"
	metaKeyAggID   MetaKey = "@@aggregate_version"
)

type Message interface {
	// ID represents the unique identity of a Message
	ID() ID
	// CausationID represents the identity of an preceding event that caused Message to happen
	CausationID() ID
	// Type contains the actual event information, or payload
	Type() interface{}
	// Recorded represents the time the event happened
	Occurred() time.Time
	// MetaData represents arbitrary data to store along the event as a hash map
	MetaData() EventMetaData
}

type Event interface {
	Message
	// AggregateID represents the identity of an aggregate
	AggregateID() ID
	// AggregateVersion represents the version of an aggregate after this event occurred
	AggregateVersion() uint64
	// WithAggregateVersion copies the Event with a replaced aggregate version
	WithAggregateVersion(uint64) Event
	// WithMetaData copies the Event with replaced EventMetaData
	WithMetaData(EventMetaData) Event
	// WithAddedMetaData copies the Event with added EventMetaData
	WithAddedMetaData(EventMetaData) Event
	// WithCausationID copies the Event and adds an causationID
	WithCausationID(ID) Event
}

// NewEvent creates a new Event
func NewEvent(aggregateID ID, payload interface{}) Event {
	return &event{
		id:          NewID(),
		eventType:   payload,
		occurred:    time.Now(),
		metaData:    EventMetaData{metaKeyAggID: aggregateID},
	}
}

func copyEvent(src *event) *event {
	return &event{
		id:               src.id,
		causationID:      src.causationID,
		eventType:        src.eventType,
		occurred:         src.occurred,
		metaData:         src.metaData.copy(),
	}
}

type eventJSON struct {
	ID               ID              `json:"id"`
	CausationID      ID              `json:"causation_id"`
	EventType        json.RawMessage `json:"type"`
	Occurred         time.Time       `json:"occurred"`
	MetaData         EventMetaData   `json:"metadata"`
}

type event struct {
	id               ID
	causationID      ID
	eventType        interface{}
	occurred         time.Time
	metaData         EventMetaData
}

func (e *event) MarshalJSON() ([]byte, error){
	evT, err := json.Marshal(e.Type())
	if err != nil {
		return []byte{}, nil
	}

	return json.Marshal(&eventJSON{
		ID:          e.ID(),
		CausationID: e.CausationID(),
		EventType:   evT,
		Occurred:    e.Occurred(),
		MetaData:    e.MetaData(),
	})
}

func (e *event) UnmarshalJSON(b []byte) error {
	j := &eventJSON{}
	if err := json.Unmarshal(b, j); err != nil {
		return err
	}

	var pl interface{}
	if err := json.Unmarshal(j.EventType, &pl); err != nil{
		return err
	}

	e.id          = j.ID
	e.causationID = j.CausationID
	e.occurred    = j.Occurred
	e.metaData    = j.MetaData
	e.eventType   = pl

	return nil
}

func (e *event) ID() ID {
	return e.id
}

func (e *event) AggregateID() ID {
	return e.metaData[metaKeyAggID].(ID)
}

func (e *event) Type() interface{} {
	return e.eventType
}

func (e *event) CausationID() ID {
	return e.causationID
}

func (e *event) Occurred() time.Time {
	return e.occurred
}

func (e *event) AggregateVersion() uint64 {
	return e.metaData[metaKeyAggVers].(uint64)
}

func (e *event) MetaData() EventMetaData {
	return e.metaData
}

func (e *event) WithAggregateVersion(v uint64) Event {
	return e.WithAddedMetaData(EventMetaData{metaKeyAggVers: v})
}

func (e *event) WithMetaData(meta EventMetaData) Event {
	cp := copyEvent(e)
	cp.metaData = meta
	return  cp
}

func (e *event) WithAddedMetaData(meta EventMetaData) Event {
	cp := copyEvent(e)
	cp.metaData = cp.metaData.merge(meta)
	return  cp
}

func (e *event) WithCausationID(cid ID) Event {
	cp := copyEvent(e)
	cp.causationID = cid
	return  cp
}
