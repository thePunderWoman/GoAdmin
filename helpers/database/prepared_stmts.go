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

func PrepareAll() error {
	PrepareAdmin()
	PrepareCurtDev()
	return nil
}

// Prepare all MySQL statements
func PrepareAdmin() error {

	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["getID"] = "select LAST_INSERT_ID() AS id"
	UnPreparedStatements["authenticateUserStmt"] = "select * from user where username=? and encpassword=? and isActive = 1"
	UnPreparedStatements["getUserByIDStmt"] = "select * from user where id=?"
	UnPreparedStatements["getUserByUsernameStmt"] = "select * from user where username=?"
	UnPreparedStatements["getUserByEmailStmt"] = "select * from user where email=?"
	UnPreparedStatements["allUserStmt"] = "select * from user"
	UnPreparedStatements["getAllModulesStmt"] = "select * from module order by module"
	UnPreparedStatements["userModulesStmt"] = "select module.* from module inner join user_module on module.id = user_module.moduleID where user_module.userID = ? order by module"
	UnPreparedStatements["setUserPasswordStmt"] = "update user set encpassword = ? where id = ?"
	UnPreparedStatements["registerUserStmt"] = "insert into user (username,email,fname,lname,isActive,superUser) VALUES (?,?,?,?,0,0)"
	UnPreparedStatements["getAllUserStmt"] = "select * from user order by fname, lname"
	UnPreparedStatements["setUserStatusStmt"] = "update user set isActive = ? WHERE id = ?"
	UnPreparedStatements["clearUserModuleStmt"] = "delete from user_module WHERE userid = ?"
	UnPreparedStatements["deleteUserStmt"] = "delete from user WHERE id = ?"
	UnPreparedStatements["addUserStmt"] = "insert into user (username,email,fname,lname,biography,photo,isActive,superUser) VALUES (?,?,?,?,?,?,?,?)"
	UnPreparedStatements["updateUserStmt"] = "update user set username=?, email=?, fname=?, lname=?, biography=?, photo=?, isActive=?, superUser=? WHERE id = ?"
	UnPreparedStatements["addModuleToUserStmt"] = "insert into user_module (userID,moduleID) VALUES (?,?)"

	if !AdminDb.IsConnected() {
		AdminDb.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareAdminStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}

	return nil
}

func PrepareCurtDev() error {
	UnPreparedStatements := make(map[string]string, 0)

	// Website Statements
	UnPreparedStatements["getAllSiteContentStmt"] = "select * from SiteContent order by page_title"
	UnPreparedStatements["getPrimaryMenuStmt"] = "select * from Menu where isPrimary = 1"
	UnPreparedStatements["getMenuItemsStmt"] = `select MSC.menuContentID, MSC.menuID, MSC.menuSort, MSC.menuTitle, MSC.menuLink, MSC.parentID, MSC.linkTarget, SC.* from Menu_SiteContent AS MSC 
												INNER JOIN Menu AS M ON MSC.menuID = M.menuID
												LEFT JOIN SiteContent AS SC ON MSC.contentID = SC.contentID
												WHERE MSC.menuID = ?`
	UnPreparedStatements["GetContentRevisionsStmt"] = "select * from SiteContentRevision WHERE contentID = ?"

	if !CurtDevDb.IsConnected() {
		CurtDevDb.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}

	return nil

}

func PrepareAdminStatement(name string, sql string, ch chan int) {
	stmt, err := AdminDb.Prepare(sql)
	if err == nil {
		Statements[name] = stmt
	} else {
		log.Println(err)
	}
	ch <- 1
}

func PrepareCurtDevStatement(name string, sql string, ch chan int) {
	stmt, err := CurtDevDb.Prepare(sql)
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
