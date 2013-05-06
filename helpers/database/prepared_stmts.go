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
	UnPreparedStatements["getAllSiteContentStmt"] = "select * from SiteContent WHERE active = 1 order by page_title"
	UnPreparedStatements["getPrimaryMenuStmt"] = "select * from Menu where isPrimary = 1"
	UnPreparedStatements["getMenuByIDStmt"] = "select * from Menu where menuID = ?"
	UnPreparedStatements["getMenuItemsStmt"] = `select MSC.menuContentID, MSC.menuID, MSC.menuSort, MSC.menuTitle, MSC.menuLink, MSC.parentID, MSC.linkTarget, SC.* from Menu_SiteContent AS MSC 
												INNER JOIN Menu AS M ON MSC.menuID = M.menuID
												LEFT JOIN SiteContent AS SC ON MSC.contentID = SC.contentID
												WHERE MSC.menuID = ?`
	UnPreparedStatements["GetContentRevisionsStmt"] = "select * from SiteContentRevision WHERE contentID = ?"
	UnPreparedStatements["GetAllMenusStmt"] = "select * from Menu where active = 1"
	UnPreparedStatements["UpdateMenuStmt"] = "Update Menu Set menu_name = ?, requireAuthentication = ?, showOnSitemap = ?, display_name = ? where menuID = ?"
	UnPreparedStatements["AddMenuStmt"] = `INSERT INTO Menu (menu_name,display_name,requireAuthentication,showOnSitemap,isPrimary,active,sort) VALUES (?,?,?,?,0,1,1)`
	UnPreparedStatements["getInsertedMenuID"] = "select LAST_INSERT_ID() FROM Menu AS id LIMIT 1"
	UnPreparedStatements["deleteMenuStmt"] = "Update Menu set active = 0 WHERE menuID = ?"
	UnPreparedStatements["clearPrimaryMenuStmt"] = "Update Menu set isPrimary = 0"
	UnPreparedStatements["setPrimaryMenuStmt"] = "Update Menu set isPrimary = 1 WHERE menuID = ?"
	UnPreparedStatements["updateMenuItemSortStmt"] = "Update Menu_SiteContent set menuSort = ?, parentID = ? WHERE menuContentID = ?"
	UnPreparedStatements["getMenuSortStmt"] = "select menuSort from Menu_SiteContent WHERE menuID = ? order by menuSort DESC"
	UnPreparedStatements["addMenuContentItemStmt"] = "INSERT INTO Menu_SiteContent (menuID,contentID,menuSort,parentID) VALUES (?,?,?,0)"
	UnPreparedStatements["addMenuLinkItemStmt"] = "INSERT INTO Menu_SiteContent (menuID,menuTitle,menuLink,linkTarget,menuSort,contentID,parentID) VALUES (?,?,?,?,?,0,0)"
	UnPreparedStatements["updateMenuLinkItemStmt"] = "UPDATE Menu_SiteContent set menuTitle = ?, menuLink = ?, linkTarget = ? WHERE menuContentID = ?"
	UnPreparedStatements["getMenuItemStmt"] = `select MSC.menuContentID, MSC.menuID, MSC.menuSort, MSC.menuTitle, MSC.menuLink, MSC.parentID, MSC.linkTarget, SC.* from Menu_SiteContent AS MSC 
												INNER JOIN Menu AS M ON MSC.menuID = M.menuID
												LEFT JOIN SiteContent AS SC ON MSC.contentID = SC.contentID
												WHERE MSC.menuContentID = ?`
	UnPreparedStatements["GetMenuParentsStmt"] = "select * from Menu_SiteContent where parentID = ? AND menuID = ? order by menuSort"
	UnPreparedStatements["DeleteMenuItemStmt"] = "delete from Menu_SiteContent where menuContentID = ?"
	UnPreparedStatements["clearPrimaryContentStmt"] = "update SiteContent set isPrimary = 0 WHERE isPrimary = 1"
	UnPreparedStatements["setPrimaryContentStmt"] = "update SiteContent set isPrimary = 1 WHERE contentID = ?"
	UnPreparedStatements["getContentStmt"] = "select * from SiteContent WHERE contentID = ?"
	UnPreparedStatements["deleteContentStmt"] = "update SiteContent set active = 0 WHERE contentID = ?"
	UnPreparedStatements["checkContentStmt"] = `select M.menu_name FROM Menu AS M
												INNER JOIN Menu_SiteContent AS MSC ON M.menuID = MSC.menuID
												WHERE MSC.contentID = ?`
	UnPreparedStatements["addContentStmt"] = `insert into SiteContent (page_title,content_type,createdDate,lastModified,meta_title,meta_description,keywords,isPrimary,published,active,slug,requireAuthentication,canonical)
											  VALUES (?,"",?,?,?,?,?,0,?,1,?,?,?)`
	UnPreparedStatements["updateContentStmt"] = `update SiteContent set page_title = ?, meta_title = ?, meta_description = ?,keywords = ?, published = ?, slug = ?, requireAuthentication = ?, canonical = ?
											     where contentID = ?`
	UnPreparedStatements["addContentRevisionStmt"] = `insert into SiteContentRevision (contentID,content_text,createdOn,active)
													  VALUES (?,?,?,?)`
	UnPreparedStatements["updateContentRevisionStmt"] = `update SiteContentRevision Set content_text = ? WHERE revisionID = ?`
	UnPreparedStatements["copyContentRevisionStmt"] = `insert into SiteContentRevision (contentID,content_text,createdOn,active)
													   (select contentID,content_text,?,0 from SiteContentRevision WHERE revisionID = ?)`
	UnPreparedStatements["getRevisionContentIDStmt"] = `select contentID from SiteContentRevision WHERE revisionID = ?`
	UnPreparedStatements["deactivateContentRevisionStmt"] = `update SiteContentRevision set active = 0 WHERE contentID = ? and active = 1`
	UnPreparedStatements["activateContentRevisionStmt"] = `update SiteContentRevision set active = 1 WHERE revisionID = ?`
	UnPreparedStatements["deleteContentRevisionStmt"] = `delete from SiteContentRevision WHERE revisionID = ?`

	// Contact Manager Statements
	UnPreparedStatements["getAllContactsStmt"] = `select * from Contact`
	UnPreparedStatements["getContactStmt"] = `select * from Contact WHERE contactID = ?`
	UnPreparedStatements["getAllContactTypesStmt"] = `select * from ContactType`
	UnPreparedStatements["getAllContactReceiversStmt"] = `select * from ContactReceiver`
	UnPreparedStatements["getContactReceiversStmt"] = `select * from ContactReceiver WHERE contactReceiverID = ?`
	UnPreparedStatements["getReceiverContactTypesStmt"] = `select CT.* from ContactType AS CT
														   INNER JOIN ContactReceiver_ContactType AS CR ON CT.contactTypeID = CR.contactTypeID
														   WHERE CR.contactReceiverID = ?`
	UnPreparedStatements["clearReceiverTypesStmt"] = `delete from ContactReceiver_ContactType WHERE contactReceiverID = ?`
	UnPreparedStatements["addReceiverTypeStmt"] = `insert into ContactReceiver_ContactType (contactReceiverID,contactTypeID) VALUES (?,?)`
	UnPreparedStatements["addContactReceiverStmt"] = `INSERT INTO ContactReceiver (first_name,last_name,email) VALUES (?,?,?)`
	UnPreparedStatements["updateContactReceiverStmt"] = `update ContactReceiver SET first_name = ?, last_name = ?, email = ? where contactReceiverID = ?`
	UnPreparedStatements["deleteContactReceiverStmt"] = `delete from ContactReceiver where contactReceiverID = ?`
	UnPreparedStatements["addContactTypeStmt"] = `insert into ContactType (name) VALUE (?)`
	UnPreparedStatements["deleteContactTypeStmt"] = `delete from ContactType WHERE contactTypeID = ?`

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
