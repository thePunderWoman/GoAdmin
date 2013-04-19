package database

import (
	"errors"
	"expvar"
	"github.com/ziutek/mymysql/mysql"
	"log"
)

// prepared statements go here
var (
	Statements = make(map[string]mysql.Stmt, 0)
)

// Prepare all MySQL statements
func PrepareAll() error {

	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["getID"] = "select LAST_INSERT_ID() AS id"
	UnPreparedStatements["authenticateUserStmt"] = "select * from user where username=? and encpassword=? and isActive = 1"
	UnPreparedStatements["getUserByIDStmt"] = "select * from user where id=?"
	UnPreparedStatements["getUserByUsernameStmt"] = "select * from user where username=?"
	UnPreparedStatements["getUserByEmailStmt"] = "select * from user where email=?"
	UnPreparedStatements["allUserStmt"] = "select * from user"
	UnPreparedStatements["userModulesStmt"] = "select module.* from module inner join user_module on module.id = user_module.moduleID where user_module.userID = ? order by module"
	UnPreparedStatements["setUserPasswordStmt"] = "update user set encpassword = ? where id = ?"
	UnPreparedStatements["registerUserStmt"] = "insert into user (username,email,fname,lname,isActive,superUser) VALUES (?,?,?,?,0,0)"
	UnPreparedStatements["getAllUserStmt"] = "select * from user order by fname, lname"

	if !AdminDb.IsConnected() {
		AdminDb.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}

	return nil
}

func PrepareStatement(name string, sql string, ch chan int) {
	stmt, err := AdminDb.Prepare(sql)
	if err == nil {
		Statements[name] = stmt
	} else {
		log.Println(err)
	}
	ch <- 1
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
