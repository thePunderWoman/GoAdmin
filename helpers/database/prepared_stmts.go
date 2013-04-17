package database

import (
	"errors"
	"expvar"
	"github.com/ziutek/mymysql/mysql"
)

// prepared statements go here
var (
	Statements map[string]mysql.Stmt
)

// Prepare all MySQL statements
func PrepareAll() error {

	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["authenticateUserStmt"] = "select * from user where username=? and encpassword=? and isActive = 1"
	UnPreparedStatements["getUserByIDStmt"] = "select * from user where id=?"
	UnPreparedStatements["getUserByUsernameStmt"] = "select * from user where username=?"
	UnPreparedStatements["getUserByEmailStmt"] = "select * from user where email=?"
	UnPreparedStatements["allUserStmt"] = "select * from user"
	UnPreparedStatements["userModulesStmt"] = "select module.* from module inner join user_module on module.id = user_module.moduleID where user_module.userID = ? order by module"
	UnPreparedStatements["setUserPasswordStmt"] = "update user set encpassword = ? where id = ?"

	Statements = make(map[string]mysql.Stmt, 0)

	if !AdminDb.IsConnected() {
		AdminDb.Connect()
	}

	for stmtname, stmtsql := range UnPreparedStatements {
		stmt, err := PrepareStatement(stmtsql)
		if err == nil {
			Statements[stmtname] = stmt
		}
	}

	return nil
}

func PrepareStatement(sql string) (stmt mysql.Stmt, err error) {
	stmt, err = AdminDb.Prepare(sql)
	return
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
