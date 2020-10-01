package identity

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewID(t *testing.T) {
	id := New()
	t.Logf("%s", id.String())
}

func TestID_Nil(t *testing.T) {
	var id ID

	if !id.Nil() {
		t.Errorf("want %#v to be nil (ish) type", id)
	}
}

func TestID_Nil_not_nil(t *testing.T) {
	id := New()
	if id.Nil() {
		t.Errorf("want %#v to be nil (ish) type", id)
	}
}

func TestID_Equals_not_equal(t *testing.T) {
	id := New()
	var idb ID

	if idb.Equals(id) {
		t.Errorf("want %s = %s", id, idb)
	}
}

func TestID_Equals(t *testing.T) {
	id := New()
	var idb ID
	copy(idb[:], id[:16])

	if !idb.Equals(id) {
		t.Errorf("want %s = %s", id, idb)
	}
}

func TestID_MarshalJSON(t *testing.T) {
	id := New()
	bt, err := id.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	if string(bt) != id.String() {
		t.Errorf("want %s got %s", bt, id)
	}
}
func TestID_UnmarshalJSON(t *testing.T) {
	var id ID
	uid, _ := uuid.NewRandom()

	if err := id.UnmarshalJSON([]byte(uid.String())); err != nil {
		t.Error(err)
	}

	if id.String() != uid.String() {
		t.Errorf("want %s got %s", uid, id)
	}
}

func TestID_Scan_stringy_bytes(t *testing.T) {
	bts := []byte(New().String())
	var id ID
	if err := id.Scan(bts); err != nil {
		t.Error(err)
	}

	if string(bts) != id.String() {
		t.Errorf("want %s, got %s", bts, id)
	}
}

func TestID_Scan_binary_bytes(t *testing.T) {
	bts := New().Bin()
	var id ID
	if err := id.Scan(bts[:]); err != nil {
		t.Error(err)
	}

	t.Logf("%s", id)

	if bts != id.Bin() {
		t.Errorf("want %#v, got %#v", bts, id)
	}
}
