package database

import (
	"errors"
	"expvar"
	"github.com/ziutek/mymysql/mysql"
)

// prepared statements go here
var (
	// example statement
	authenticateUserStmt  = "select * from user where username=? and encpassword=? and isActive = 1"
	getUserByIDStmt       = "select * from user where id=?"
	getUserByUsernameStmt = "select * from user where username=?"
	getUserByEmailStmt    = "select * from user where email=?"
	allUserStmt           = "select * from user"
	userModulesStmt       = "select module.* from module inner join user_module on module.id = user_module.moduleID where user_module.userID = ? order by module"
	setUserPasswordStmt   = "update user set encpassword = ? where id = ?"
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

	getUserByUsernamePrepared, err := AdminDb.Prepare(getUserByUsernameStmt)
	if err != nil {
		return err
	}
	Statements["getUserByUsernameStmt"] = getUserByUsernamePrepared

	getUserByEmailPrepared, err := AdminDb.Prepare(getUserByEmailStmt)
	if err != nil {
		return err
	}
	Statements["getUserByEmailStmt"] = getUserByEmailPrepared

	userModulesPrepared, err := AdminDb.Prepare(userModulesStmt)
	if err != nil {
		return err
	}
	Statements["userModulesStmt"] = userModulesPrepared

	setUserPasswordPrepared, err := AdminDb.Prepare(setUserPasswordStmt)
	if err != nil {
		return err
	}
	Statements["setUserPasswordStmt"] = setUserPasswordPrepared

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
