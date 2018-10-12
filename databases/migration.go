package databases


// make database structure
type Migrator interface {
	migrate()
}