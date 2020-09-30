package message

import (
	"time"
)

type Event interface {
	ID() ID
	CausationID() ID
	AggregateID() ID
	AggregateVersion() uint64
	Type() interface{}
	Recorded() time.Time
	MetaData() EventMetaData
	WithAggregateVersion(uint64) Event
	WithMetaData(EventMetaData) Event
	WithAddedMetaData(EventMetaData) Event
	WithCausationID(ID) Event
}

type event struct {
	id ID
	causationID ID
	aggregateID ID
	aggregateVersion uint64
	eventType interface{}
	recorded time.Time
	metaData EventMetaData
}

func (e *event) ID() ID {
	return e.id
}

func (e *event) AggregateID() ID {
	return e.aggregateID
}

func (e *event) Type() interface{} {
	return e.eventType
}

func (e *event) CausationID() ID {
	return e.causationID
}

func (e *event) Recorded() time.Time {
	return e.recorded
}

func (e *event) AggregateVersion() uint64 {
	return e.aggregateVersion
}

func (e *event) MetaData() EventMetaData {
	return e.metaData
}

func (e *event) WithAggregateVersion(v uint64) Event {
	cp := copyEvent(e)
	cp.aggregateVersion = v
	return  cp
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

func copyEvent(src *event) *event {
	return &event{
		id:               src.id,
		causationID:      src.causationID,
		aggregateID:      src.aggregateID,
		aggregateVersion: src.aggregateVersion,
		eventType:        src.eventType,
		recorded:         src.recorded,
		metaData:         src.metaData,
	}
}


func NewEvent(aid ID, eType interface{}) Event {
	return &event{
		id:               NewID(),
		aggregateID:      aid,
		eventType:        eType,
		recorded:         time.Now(),
		metaData:         make(EventMetaData),
	}
}
