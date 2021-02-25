package aggregate

type Snapshot interface {
	Version() uint64
	State() interface{}
}

type SnapShooter interface {
	ApplySnapshot(Snapshot) error
	Snapshot() (Snapshot, error)
}
