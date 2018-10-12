package databases

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

func Connection(database TypeDatabase)  {
	db, err := sql.Open(database.String(), "mysql://root:1234@localhost/payout")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//
	//_,err = db.Exec("CREATE DATABASE "+name)
	//if err != nil {
	//	panic(err)
	//}
	//
	//_,err = db.Exec("USE "+name)
	//if err != nil {
	//	panic(err)
	//}
	//
	//_,err = db.Exec("CREATE TABLE example ( id integer, data varchar(32) )")
	//if err != nil {
	//	panic(err)
	//}
}