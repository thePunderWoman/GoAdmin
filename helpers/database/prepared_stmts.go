package database

import (
	"errors"
	"expvar"
	"github.com/ziutek/mymysql/mysql"
)

// prepared statements go here
var (
	// example statement
	authenticateUserStmt = "select * from user where username=? and encpassword=? and isActive = 1"
	getUserByIDStmt      = "select * from user where id=?"
	allUserStmt          = "select * from user"
	userModulesStmt      = "select module.* from module inner join user_module on module.id = user_module.moduleID where user_module.userID = ? order by module"
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

	allUserPrepared, err := AdminDb.Prepare(allUserStmt)
	if err != nil {
		return err
	}
	Statements["allUserStmt"] = allUserPrepared

	getUserByIDPrepared, err := AdminDb.Prepare(getUserByIDStmt)
	if err != nil {
		return err
	}
	Statements["getUserByIDStmt"] = getUserByIDPrepared

	userModulesPrepared, err := AdminDb.Prepare(userModulesStmt)
	if err != nil {
		return err
	}
	Statements["userModulesStmt"] = userModulesPrepared

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
