package databases

type NameDatabase interface {
	GetType() string
}

// make database structure
type Migrator interface {
	Migrate() string
}

type DriverDatabase interface {
	NameDatabase
	Migrator
}