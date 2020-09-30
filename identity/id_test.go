package identity

import "testing"

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

