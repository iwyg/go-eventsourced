package identity

import (
	"github.com/google/uuid"
)


type ID [16]byte
var NIL ID

func (id ID) Bin() [16]byte {
	return id
}

func (id ID) Nil() bool {
	return id == NIL
}

func (id ID) String() string {
	u, _ := uuid.FromBytes(id[:])
	return u.String()
}

func (id ID) Equals(other ID) bool {
	return id == other
}

func (id ID) Validate() error {
	_, err := uuid.FromBytes(id[:])
	return err
}

func New() ID {
	var id ID
	uid, _ := uuid.NewRandom()
	idb, _ := uid.MarshalBinary()
	copy(id[:], idb[:16])
	return id
}
