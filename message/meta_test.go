package message

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestEventMetaData_Scan(t *testing.T) {
	var md MetaData
	src := []byte(`{"answerToEverything":42}`)

	if err := md.Scan(src); err != nil {
		t.Error(err)
	}

	a, has := md["answerToEverything"]
	if !has {
		t.Errorf("key %q not found in MetaData", "answerToEverything")
	}

	if a != float64(42) {
		t.Errorf("want %q, got %q", 42, a)
	}
}

func TestEventMetaData_MarshalJSON_with_id(t *testing.T) {
	id := NewID()
	md := MetaData{metaKeyAggID: id}
	b, err := json.Marshal(md)
	if err != nil {
		t.Error(err)
	}
	v := fmt.Sprintf("{%q:%q}", metaKeyAggID, id.String())

	if string(b) != v {
		t.Errorf("want %s, got %s", v, string(b))
	}
}

func TestEventMetaData_Scan_Reconstitute_ID(t *testing.T) {
	var md MetaData
	src := []byte(`{"@@aggregate_id":"25410a5a-24bf-11eb-bcb1-d362a6499ff0"}`)

	if err := json.Unmarshal(src, &md); err != nil {
		t.Error(err)
	}

	v := md[metaKeyAggID]
	switch v.(type) {
	case ID:
		return
	default:
		t.Errorf("want %s, got %s", "message.ID", reflect.TypeOf(v).Name())
	}
}
