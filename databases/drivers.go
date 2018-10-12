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
	RequestCommit(conn *sql.DB)
	ResponseCommit(conn *sql.DB)
}

type RequestTable struct {
	type_transfer string
	dstAccount string
	clientOrderId uint64
	requestDT string
	amount float64
	currency uint32
	agentId uint64
	contract string
	paymentParams string
}
