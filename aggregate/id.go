package aggregate

import "github.com/iwyg/go-eventsourced/identity"

type ID identity.ID

func NewID() ID {
	return ID(identity.New())
}

func (id ID) Equals(other ID) bool {
	return identity.ID(id).Equals(identity.ID(other))
}

func (id ID) String() string {
	return identity.ID(id).String()
}

func (id ID) Nil() bool {
	return identity.ID(id).Nil()
}
