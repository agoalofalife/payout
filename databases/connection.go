package databases

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func Connection(databaseDriver DriverDatabase, login string, password string, host string, table string)  {
	db, err := sql.Open(databaseDriver.GetType(), fmt.Sprintf("%s:%s@%s/%s",login, password, host, table))

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		log.Println("Database is success connected!")
	}
	// migrate structure
	_,err = db.Exec(databaseDriver.Migrate())
	if err != nil {
		panic(err)
	} else {
		log.Println("Success migrate structure table!")
	}

}