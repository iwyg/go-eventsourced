package aggregate

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"testing"
)

func TestID_Equals(t *testing.T) {
	var aggCopy ID
	aggID := NewID()
	if err := aggCopy.Scan([]byte(aggID.String())); err != nil {
		t.Error(err)
	}

	if !aggID.Equals(aggCopy) {
		t.Errorf("ID %#v should equal %#v", aggID, aggCopy)
	}
}

func TestID_Scan(t *testing.T) {
	var aggCopy ID
	uid, _ := uuid.NewRandom()
	if err := aggCopy.Scan([]byte(uid.String())); err != nil {
		t.Error(err)
	}
}

func TestID_UnmarshalJSON(t *testing.T) {
	var aggID ID
	uid, _ := uuid.NewRandom()
	if err := json.Unmarshal([]byte(fmt.Sprintf("%q", uid.String())), &aggID); err != nil {
		t.Error(err)
	}

	if uid.String() != aggID.String() {
		t.Errorf("want %s, got %s", uid.String(), aggID.String())
	}

}

func TestID_MarshalJSON(t *testing.T) {
	aggID := NewID()

	bts, err :=  json.Marshal(aggID)
	if err != nil {
		t.Error(err)
	}

	raw, _ :=  strconv.Unquote(string(bts))

	if raw != aggID.String() {
		t.Errorf("want %q, got %q", aggID.String(), raw)
	}

}
