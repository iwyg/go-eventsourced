package message

import (
	"testing"
	"time"
)

func TestEvent_Occurred(t *testing.T) {
	id := NewID()
	now := time.Now()
	ev := NewEvent(id, nil)

	if now.Unix() != ev.Occurred().Unix() {
		t.Errorf("should be equal")
	}
}

func TestEvent_WithMetaData(t *testing.T) {
	id   := NewID()
	ev   := NewEvent(id, struct {}{})
	meta := MetaData{"__user": "username"}
	ev    = ev.WithMetaData(meta)
	eventCopy := ev.WithMetaData(MetaData{"__roles": []string{"admin", "user"}})

	if _, found := ev.MetaData()["__user"]; !found {
		t.Errorf("key %q should not in original meta data", "__user")
	}

	if _, found := eventCopy.MetaData()["__user"]; found {
		t.Errorf("key %q should not be in meta data", "__user")
	}

	if _, found := eventCopy.MetaData()["__roles"]; !found {
		t.Errorf("key %q should be in meta data", "__roles")
	}
}

func TestEvent_WithAddedMetaData(t *testing.T) {
	id   := NewID()
	ev   := NewEvent(id, struct {}{})
	meta := MetaData{"__user": "username"}
	ev    = ev.WithMetaData(meta)
	eventCopy := ev.WithAddedMetaData(MetaData{"__roles": []string{"admin", "user"}})

	if _, found := ev.MetaData()["__user"]; !found {
		t.Errorf("key %q should be in original meta data", "__user")
	}

	if _, found := ev.MetaData()["__roles"]; found {
		t.Errorf("key %q should not be in original meta data", "__roles")
	}

	if _, found := eventCopy.MetaData()["__user"]; !found {
		t.Errorf("key %q should be in added meta data", "__user")
	}

	if _, found := eventCopy.MetaData()["__roles"]; !found {
		t.Errorf("key %q should be in added meta data", "__roles")
	}
}
