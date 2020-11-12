package message

import (
	"encoding/json"
	"errors"
)

type MetaKey string
func (mk MetaKey) String() string {
	return string(mk)
}

type MetaData map[MetaKey]interface{}

func (em *MetaData) Scan(src interface{}) error {
	p, ok := src.([]byte)
	if !ok {
		return errors.New("can't cast source input to []byte")
	}

	return json.Unmarshal(p, em)
}

func (em MetaData) merge(meta MetaData) MetaData {
	m := em
	if meta == nil {
		return m
	}

	for key, val := range meta {
		m[key] = val
	}

	return m
}
func (em MetaData) copy() MetaData {
	m := make(MetaData, len(em))
	for key, val := range em {
		m[key] = val
	}

	return m
}
