package message

import (
	"encoding/json"
	"errors"
)

type MetaKey string
func (mk MetaKey) String() string {
	return string(mk)
}

type EventMetaData map[MetaKey]interface{}

func (em *EventMetaData) Scan(src interface{}) error {
	p, ok := src.([]byte)
	if !ok {
		return errors.New("can't cast source input to []byte")
	}

	return json.Unmarshal(p, em)
}

func (em EventMetaData) merge(meta EventMetaData) EventMetaData {
	m := em
	if meta == nil {
		return m
	}

	for key, val := range meta {
		m[key] = val
	}

	return m
}
func (em EventMetaData) copy() EventMetaData {
	m := make(EventMetaData, len(em))
	for key, val := range em {
		m[key] = val
	}

	return m
}
