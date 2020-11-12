package message

import "testing"

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
