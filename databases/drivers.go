package databases

import "database/sql"

type NameDatabase interface {
	GetType() string
}

// make database structure
type Migrator interface {
	Migrate() []string
}

type DriverDatabase interface {
	NameDatabase
	Migrator
}

type Commiter interface {
	requestCommit(conn *sql.DB)
	responseCommit(conn *sql.DB)
}
