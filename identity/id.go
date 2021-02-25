package identity

import (
	"github.com/google/uuid"
	"strconv"
)


type ID [16]byte
var NIL ID

func (id ID) Bin() [16]byte {
	return id
}

func (id ID) Nil() bool {
	return id == NIL
}

func (id ID) MarshalJSON() ([]byte, error)  {
	return []byte(strconv.Quote(id.String())), nil
}

func (id *ID) UnmarshalJSON(b []byte) error  {
	uid, err := uuid.ParseBytes(b)
	if err != nil {
		return err
	}
	idb, err := uid.MarshalBinary()
	if err != nil {
		return err
	}

	copy(id[:], idb[:16])
	return nil
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func (id *ID) Scan(src interface{}) error {
	uid := uuid.UUID(*id)
	if err := uid.Scan(src); err != nil {
		return err
	}

	copy(id[:], uid[:])

	return nil
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
