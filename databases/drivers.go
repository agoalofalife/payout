package databases

import "database/sql"

const (
	TransferPhone = iota
	TransferPurse
)
type TypeTransfer int
var typesTransfer = [...]string {"phone", "purse",}
func (transfer TypeTransfer) String() {
	return
}

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
	TypeTransfer
	DstAccount string
	ClientOrderId uint64
	RequestDT string
	Amount float64
	Currency uint32
	AgentId uint64
	Contract string
	PaymentParams string
}
