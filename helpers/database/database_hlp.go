package database

import (
	"github.com/ziutek/mymysql/thrsafe"
	"log"
	"os"
)

var (
	// MySQL Connection Handler
	CurtDevDb = thrsafe.New(db_proto, "", db_addr, db_user, db_pass, CurtDevdb_name)
	AdminDb   = thrsafe.New(db_proto, "", db_addr, db_user, db_pass, Admindb_name)
)

func MysqlError(err error) (ret bool) {
	ret = (err != nil)
	if ret {
		log.Println("MySQL error: ", err)
	}
	return
}

func MysqlErrExit(err error) {
	if MysqlError(err) {
		os.Exit(1)
	}
}
