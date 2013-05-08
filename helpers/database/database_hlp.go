package database

import (
	"github.com/ziutek/mymysql/autorc"
	"log"
	"os"
)

var (
	// MySQL Connection Handler
	CurtDevDb = autorc.New(db_proto, "", db_addr, db_user, db_pass, CurtDevdb_name)
	AdminDb   = autorc.New(db_proto, "", db_addr, db_user, db_pass, Admindb_name)
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
