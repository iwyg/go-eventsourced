package message

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (em *MetaData) UnmarshalJSON(src []byte) error {
	var tg map[MetaKey]interface{}
	if err := json.Unmarshal(src, &tg); err != nil {
		return nil
	}
	for _, c := range []MetaKey{metaKeyAggID} {
		if v, f := tg[c]; f {
			b, err := stringLikeToBytes(v)
			if err != nil {
				return err
			}

			var id ID
			if err := id.Scan(b); err != nil {
				return err
			}
			tg[c] = id
		}
	}
	*em = tg
	return nil
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

func stringLikeToBytes(v interface{}) ([]byte, error) {
	switch v.(type) {
	case string:
		s := v.(string)
		return []byte(s), nil
	case []byte:
		return v.([]byte), nil
	default:
		return nil, errors.New(fmt.Sprintf("cant cast %#v to []byte", v))
	}
}