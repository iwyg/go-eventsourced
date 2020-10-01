package aggregate

import "github.com/iwyg/go-eventsourced/identity"

type ID identity.ID

func NewID() ID {
	return ID(identity.New())
}

func (id ID) Equals(other ID) bool {
	return id.id().Equals(other.id())
}

func (id ID) MarshalJSON() ([]byte, error) {
	return id.id().MarshalJSON()
}

func (id ID) UnmarshalJSON(b []byte) error {
	recv := id.id()
	return recv.UnmarshalJSON(b)
}

func (id ID) String() string {
	return id.id().String()
}

func (id *ID) Scan(src interface{}) error {
	var n identity.ID
	if err := n.Scan(src); err != nil {
		return err
	}

	copy(id[:], n[:])
	return nil
}

func (id ID) Nil() bool {
	return id.id().Nil()
}

func (id ID) id() identity.ID {
	return identity.ID(id)
}