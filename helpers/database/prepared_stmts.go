package database

import (
	"errors"
	"expvar"
	"github.com/ziutek/mymysql/mysql"
)

// prepared statements go here
var (
	// example statement
	authenticateUserStmt = "select * from user where username=? and password=?"
)

// Create map of all statements
var (
	Statements map[string]mysql.Stmt
)

// Prepare all MySQL statements
func PrepareAll() error {

	Statements = make(map[string]mysql.Stmt, 0)

	if !AdminDb.IsConnected() {
		AdminDb.Connect()
	}

	// Example Preparation
	authenticateUserPrepared, err := AdminDb.Prepare(authenticateUserStmt)
	if err != nil {
		return err
	}
	Statements["authenticateUserStmt"] = authenticateUserPrepared

	return nil
}

func GetStatement(key string) (stmt mysql.Stmt, err error) {
	stmt, ok := Statements[key]
	if !ok {
		qry := expvar.Get(key)
		if qry == nil {
			err = errors.New("Invalid query reference")
		}
	}
	return

}
