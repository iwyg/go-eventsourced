package message

import (
	"encoding/json"
	"testing"
)

func TestID_Scan_stringy_bytes(t *testing.T) {
	bts := []byte(NewID().String())
	var id ID
	if err := id.Scan(bts); err != nil {
		t.Error(err)
	}

	if string(bts) != id.String() {
		t.Errorf("want %s, got %s", bts, id)
	}
}

func TestID_Scan_binary_bytes(t *testing.T) {
	bts := NewID().id().Bin()
	var id ID
	if err := id.Scan(bts[:]); err != nil {
		t.Error(err)
	}

	t.Logf("%s", id)

	if bts != id.id().Bin() {
		t.Errorf("want %#v, got %#v", bts, id)
	}
}

func TestID_MarshalJSON(t *testing.T) {
	guid := NewID()
	t.Logf("%s", guid)
	b, err :=json.Marshal(guid)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v", b)
}
