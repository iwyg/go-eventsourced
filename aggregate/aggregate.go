package aggregate

import "github.com/iwyg/go-eventsourced/identity"

type Aggregate interface {
	// ID return the aggregate identity
	ID() identity.ID
	// Version return the aggregate version
	Version() uint64
}
