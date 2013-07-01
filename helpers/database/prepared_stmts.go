package database

import (
	"errors"
	"expvar"
	"github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"log"
)

// prepared statements go here
var (
	Statements = make(map[string]*autorc.Stmt, 0)
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

	if !AdminDb.Raw.IsConnected() {
		AdminDb.Raw.Connect()
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

	// FAQ Manager Statements
	UnPreparedStatements["GetAllFAQStmt"] = `select * from FAQ`
	UnPreparedStatements["GetFAQStmt"] = `select * from FAQ WHERE faqID = ?`
	UnPreparedStatements["UpdateFAQStmt"] = `UPDATE FAQ SET question = ?, answer = ? WHERE faqID = ?`
	UnPreparedStatements["AddFAQStmt"] = `INSERT INTO FAQ (question, answer) VALUES (?,?)`
	UnPreparedStatements["DeleteFAQStmt"] = `DELETE FROM FAQ WHERE faqID = ?`

	// News Manager Statements
	UnPreparedStatements["GetAllNewsStmt"] = `select * from NewsItem WHERE active = 1`
	UnPreparedStatements["GetNewsItemStmt"] = `select * from NewsItem WHERE newsItemID = ?`
	UnPreparedStatements["UpdateNewsItemStmt"] = `Update NewsItem SET title = ?, lead = ?, content = ?, publishStart = ?, publishEnd = ?, slug = ? WHERE newsItemID = ?`
	UnPreparedStatements["AddNewsItemStmt"] = `INSERT INTO NewsItem (title,lead,content,publishStart,publishEnd,slug,active) VALUES (?,?,?,?,?,?,1)`
	UnPreparedStatements["DeleteNewsItemStmt"] = `Update NewsItem SET active = 0 WHERE newsItemID = ?`

	// Video Manager Statements
	UnPreparedStatements["GetAllVideosStmt"] = `select * from Video`
	UnPreparedStatements["GetVideoStmt"] = `select * from Video where videoID = ?`
	UnPreparedStatements["UpdateVideoSortStmt"] = `Update Video Set sort = ? where videoID = ?`
	UnPreparedStatements["DeleteVideoStmt"] = `DELETE FROM Video where videoID = ?`
	UnPreparedStatements["GetLastVideoSortStmt"] = `select sort from Video Order By sort desc`
	UnPreparedStatements["AddVideoStmt"] = `INSERT INTO Video (embed_link,dateAdded,sort,title,description,youtubeID,watchpage,screenshot) VALUES (?,?,?,?,?,?,?,?)`

	//Testimonial Manager Statements
	UnPreparedStatements["GetAllTestimonialsStmt"] = `SELECT * from Testimonial WHERE active = 1 AND approved = ?`
	UnPreparedStatements["DeleteTestimonialStmt"] = `UPDATE Testimonial SET active = 0 WHERE testimonialID = ?`
	UnPreparedStatements["SetTestimonialApprovalStmt"] = `UPDATE Testimonial SET approved = ? WHERE testimonialID = ?`

	//Landing Page Manager Statements
	UnPreparedStatements["GetActiveLandingPagesStmt"] = `SELECT * from LandingPage WHERE endDate > UTC_TIMESTAMP()`
	UnPreparedStatements["GetPastLandingPagesStmt"] = `SELECT * from LandingPage WHERE endDate < UTC_TIMESTAMP()`
	UnPreparedStatements["GetLandingPageImagesStmt"] = `SELECT * from LandingPageImages WHERE landingPageID = ?`
	UnPreparedStatements["GetLandingPageDataStmt"] = `SELECT * from LandingPageData WHERE landingPageID = ?`
	UnPreparedStatements["AddLandingPageStmt"] = `INSERT INTO LandingPage (name,startDate,endDate,url,newWindow,menuPosition) VALUES (?,?,?,?,?,?)`
	UnPreparedStatements["UpdateLandingPageStmt"] = `UPDATE LandingPage SET name = ?, startDate = ?, endDate = ?, url = ?, pageContent = ?, linkClasses = ?, conversionID = ?, conversionLabel = ?, newWindow = ?, menuPosition = ? WHERE id = ?`
	UnPreparedStatements["GetLandingPageStmt"] = `SELECT * from LandingPage WHERE id = ?`
	UnPreparedStatements["AddLandingPageDataStmt"] = `INSERT INTO LandingPageData (landingPageID,dataKey,dataValue) VALUES (?,?,?)`
	UnPreparedStatements["DeleteLandingPageDataStmt"] = `DELETE FROM LandingPageData WHERE id = ?`
	UnPreparedStatements["DeleteLandingPageStmt"] = `DELETE FROM LandingPage WHERE id = ?`

	// Sales Reps
	UnPreparedStatements["GetAllSalesRepsStmt"] = `SELECT salesRepID, name, code, (select COUNT(cust_id) from Customer WHERE Customer.salesRepID = SalesRepresentative.salesRepID) AS customercount from SalesRepresentative`
	UnPreparedStatements["GetSalesRepStmt"] = `SELECT salesRepID, name, code, (select COUNT(cust_id) from Customer WHERE Customer.salesRepID = SalesRepresentative.salesRepID) AS customercount from SalesRepresentative WHERE salesRepID = ?`
	UnPreparedStatements["UpdateSalesRepStmt"] = `UPDATE SalesRepresentative set name = ?, code = ? WHERE salesRepID = ?`
	UnPreparedStatements["AddSalesRepStmt"] = `INSERT INTO SalesRepresentative (name,code) VALUES (?,?)`
	UnPreparedStatements["DeleteSalesRepStmt"] = `DELETE FROM SalesRepresentative WHERE salesRepID = ?`

	// Blog
	UnPreparedStatements["GetAllPostsStmt"] = `SELECT * from BlogPosts WHERE active = 1`
	UnPreparedStatements["GetPostStmt"] = `SELECT * from BlogPosts WHERE blogPostID = ?`
	UnPreparedStatements["GetAllBlogCategoriesStmt"] = `SELECT * from BlogCategories WHERE active = 1`
	UnPreparedStatements["GetBlogCategoryStmt"] = `SELECT * from BlogCategories WHERE blogCategoryID = ?`
	UnPreparedStatements["AddBlogCategoryStmt"] = `INSERT INTO BlogCategories (name,slug,active) VALUES (?,?,1)`
	UnPreparedStatements["UpdateBlogCategoryStmt"] = `Update BlogCategories set name = ?, slug = ? WHERE blogCategoryID = ?`
	UnPreparedStatements["DeleteBlogCategoryStmt"] = `Update BlogCategories set active = 0 WHERE blogCategoryID = ?`
	UnPreparedStatements["GetPostCategoriesStmt"] = `SELECT BC.* from BlogCategories AS BC
													 INNER JOIN BlogPost_BlogCategory AS BPBC ON BC.blogCategoryID = BPBC.blogCategoryID
													 WHERE blogPostID = ? AND BC.active = 1`
	UnPreparedStatements["GetPostCommentsStmt"] = `SELECT * from Comments WHERE blogPostID = ? AND active = 1`
	UnPreparedStatements["AddPostStmt"] = `INSERT INTO BlogPosts (post_title,slug,post_text,createdDate,userID,meta_title,meta_description,keywords,active) VALUES (?,?,?,UTC_TIMESTAMP(),?,?,?,?,1)`
	UnPreparedStatements["UpdatePostStmt"] = `UPDATE BlogPosts SET post_title = ?, slug = ?, post_text = ?, userID = ?, meta_title = ?, meta_description = ?, keywords = ? WHERE blogPostID = ?`
	UnPreparedStatements["PublishPostStmt"] = `UPDATE BlogPosts SET publishedDate = UTC_TIMESTAMP() WHERE blogPostID = ?`
	UnPreparedStatements["UnPublishPostStmt"] = `UPDATE BlogPosts SET publishedDate = null WHERE blogPostID = ?`
	UnPreparedStatements["ClearPostCategoriesStmt"] = `DELETE FROM BlogPost_BlogCategory WHERE blogPostID = ?`
	UnPreparedStatements["AddPostCategoryStmt"] = `INSERT INTO BlogPost_BlogCategory (blogPostID,blogCategoryID) VALUES (?,?)`
	UnPreparedStatements["DeletePostStmt"] = `UPDATE BlogPosts SET active = 0 WHERE blogPostID = ?`
	UnPreparedStatements["GetBlogCommentsStmt"] = `SELECT * from Comments WHERE active = 1 AND approved = 0`
	UnPreparedStatements["GetBlogCommentStmt"] = `SELECT * from Comments WHERE commentID = ?`
	UnPreparedStatements["ApproveBlogCommentStmt"] = `UPDATE Comments set approved = 1 WHERE commentID = ?`
	UnPreparedStatements["DeleteBlogCommentStmt"] = `UPDATE Comments set active = 0 WHERE commentID = ?`

	// Locations
	UnPreparedStatements["GetAllCountriesStmt"] = `SELECT * from Country`
	UnPreparedStatements["GetCountryStmt"] = `SELECT * from Country WHERE countryID = ?`
	UnPreparedStatements["GetStatesByCountryStmt"] = `SELECT * from States WHERE countryID = ?`
	UnPreparedStatements["GetAllStatesStmt"] = `SELECT * from States`
	UnPreparedStatements["GetStateStmt"] = `SELECT * from States WHERE stateID = ?`

	// Customer
	UnPreparedStatements["GetAllCustomersStmt"] = `SELECT *, 
													(SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id) AS locationCount,
													(SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id) AS userCount,
													(SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id) AS propertyCount 
													from Customer`
	UnPreparedStatements["GetAllSimpleCustomersStmt"] = `SELECT cust_id, name, customerID from Customer`
	UnPreparedStatements["GetCustomerStmt"] = `SELECT *, 
												(SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id) AS locationCount,
												(SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id) AS userCount,
												(SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id) AS propertyCount 
												from Customer WHERE cust_id = ?`
	UnPreparedStatements["UpdateCustomerStmt"] = `UPDATE Customer SET name = ?, email = ?, address = ?, address2 = ?, city = ?, stateID = ?, postal_code = ?, phone = ?, fax = ?, 
												  contact_person = ?, dealer_type = ?, tier = ?, website = ?, searchURL = ?, eLocalURL = ?, logo = ?, customerID = ?,
												  parentID = ?, isDummy = ?, mCodeID = ?, salesRepID = ?, showWebsite = ? WHERE cust_id = ?`
	UnPreparedStatements["AddCustomerStmt"] = `INSERT INTO Customer (name,email,address,address2,city,stateID,postal_code,phone,fax,contact_person,dealer_type,tier,website,searchURL,eLocalURL,logo,customerID,parentID,isDummy,mCodeID,salesRepID,showWebsite)
											   VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	// Dealer Types
	UnPreparedStatements["GetAllDealerTypesStmt"] = `SELECT * From DealerTypes`
	UnPreparedStatements["GetDealerTypeStmt"] = `SELECT * from DealerTypes WHERE dealer_type = ?`

	// Dealer Tier
	UnPreparedStatements["GetAllDealerTiersStmt"] = `SELECT * from DealerTiers`
	UnPreparedStatements["GetDealerTierStmt"] = `SELECT * from DealerTiers WHERE ID = ?`

	// Mapics Codes
	UnPreparedStatements["GetAllMapicsCodesStmt"] = `SELECT * from MapixCode`
	UnPreparedStatements["GetMapicsCodeStmt"] = `SELECT * from MapixCode WHERE mCodeID = ?`

	// MapIcons
	UnPreparedStatements["GetMapIconsStmt"] = `SELECT * from MapIcons`

	// Customer Locations
	UnPreparedStatements["GetCustomerLocationsStmt"] = `SELECT * from CustomerLocations WHERE cust_id = ?`
	UnPreparedStatements["GetCustomerLocationsNoGeoStmt"] = `SELECT * from CustomerLocations WHERE cust_id = ? AND (latitude = 0 OR longitude = 0)`
	UnPreparedStatements["GetCustomerLocationStmt"] = `SELECT * from CustomerLocations WHERE locationID = ?`
	UnPreparedStatements["UpdateCustomerLocationStmt"] = `UPDATE CustomerLocations SET name=?, address = ?, city = ?, stateID = ?, postalCode = ?, email = ?, phone = ?, fax = ?, latitude = ?, longitude = ? WHERE locationID = ?`
	UnPreparedStatements["AddCustomerLocationStmt"] = `INSERT INTO CustomerLocations (name,address,city,stateID,postalCode,email,phone,fax,latitude,longitude,cust_id,isprimary,ShippingDefault,contact_person) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	UnPreparedStatements["DeleteCustomerLocationStmt"] = `DELETE FROM CustomerLocations WHERE locationID = ?`

	// Customer Users
	UnPreparedStatements["GetAllCustomerUsersStmt"] = `SELECT * from CustomerUser`
	UnPreparedStatements["GetCustomerUsersStmt"] = `SELECT * from CustomerUser WHERE cust_id = ?`
	UnPreparedStatements["GetCustomerUserKeysStmt"] = `select AK.*, AKT.type, AKT.date_added AS typeDateAdded
														from ApiKey AK
														INNER JOIN ApiKeyType AKT ON AK.type_id = AKT.id
														INNER JOIN CustomerUser CU on AK.user_id = CU.id
														where CU.cust_ID = ?`
	UnPreparedStatements["GetAllCustomerUserKeysStmt"] = `select AK.*, AKT.type, AKT.date_added AS typeDateAdded
														from ApiKey AK
														INNER JOIN ApiKeyType AKT ON AK.type_id = AKT.id`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
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

func GetStatement(key string) (stmt *autorc.Stmt, err error) {
	stmt, ok := Statements[key]
	if !ok {
		qry := expvar.Get(key)
		if qry == nil {
			err = errors.New("Invalid query reference")
		}
	}
	return

}
