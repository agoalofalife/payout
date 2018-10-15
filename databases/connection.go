package databases

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"sync"
)

var (
	db *sql.DB
	once sync.Once
)
func Connection(databaseDriver DriverDatabase, login string, password string, host string, table string) *sql.DB {
	var err error

	once.Do(func() {
		db, err = sql.Open(databaseDriver.GetType(), fmt.Sprintf("%s:%s@%s/%s",login, password, host, table))
		if err != nil {
			panic(err)
		}

		err = db.Ping()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			log.Println("Database is success connected!")
		}
		// migrate structure
		err = createTable(db, databaseDriver.Migrate())
		if err != nil {
			panic(err)
		} else {
			log.Println("Success migrate structure table!")
		}
	})
	//defer db.Close()

	return  db
}


// createTable creates the table, and if necessary, the database.
func createTable(conn *sql.DB, createTableStatements []string) error {
	for _, stmt := range createTableStatements {
		_, err := conn.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
