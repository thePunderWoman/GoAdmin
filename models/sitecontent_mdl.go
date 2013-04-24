package models

import (
	"../helpers/database"
	_ "errors"
	"github.com/ziutek/mymysql/mysql"
	_ "log"
	"time"
)

type Menu struct {
	ID            int
	Name          string
	Primary       bool
	Active        bool
	DisplayName   string
	RequireAuth   bool
	ShowOnSitemap bool
	Sort          int
	Items         []MenuItem
}

type MenuItem struct {
	ID         int
	MenuID     int
	ContentID  int
	Sort       int
	Title      string
	Link       string
	ParentID   int
	LinkTarget string
	Content    Content
}

type Content struct {
	ID              int
	ContentType     string
	PageTitle       string
	CreatedDate     time.Time
	LastModified    time.Time
	MetaTitle       string
	MetaDescription string
	Keywords        string
	Primary         bool
	Published       bool
	Active          bool
	Slug            string
	RequireAuth     bool
	Canonical       string
	Revisions       []ContentRevision
	ActiveRevision  ContentRevision
}

type ContentRevision struct {
	ID          int
	ContentID   int
	ContentText string
	CreatedOn   time.Time
	Active      bool
}

func GetAllSiteContent() (contents []Content, err error) {
	sel, err := database.GetStatement("getAllSiteContentStmt")
	if err != nil {
		return contents, err
	}
	rows, res, err := sel.Exec()
	if database.MysqlError(err) {
		return contents, err
	}

	c := make(chan Content)

	for _, row := range rows {
		go PopulateContent(row, res, c)
	}

	for _, _ = range rows {
		contents = append(contents, <-c)
	}

	return contents, nil
}

func PopulateContent(row mysql.Row, res mysql.Result, ch chan Content) {
	cid := res.Map("contentID")
	contentType := res.Map("content_type")
	pageTitle := res.Map("page_title")
	createdDate := res.Map("createdDate")
	lastModified := res.Map("lastModified")
	metaTitle := res.Map("meta_title")
	metaDesc := res.Map("meta_description")
	keywords := res.Map("keywords")
	isPrimary := res.Map("isPrimary")
	published := res.Map("published")
	active := res.Map("active")
	slug := res.Map("slug")
	requireAuth := res.Map("requireAuthentication")
	canonical := res.Map("canonical")
	var content Content

	id := row.Int(cid)
	revCh := make(chan []ContentRevision)
	go GetContentRevisions(id, revCh)

	content = Content{
		ID:              id,
		ContentType:     row.Str(contentType),
		PageTitle:       row.Str(pageTitle),
		CreatedDate:     row.Time(createdDate, UTC),
		LastModified:    row.Time(lastModified, UTC),
		MetaTitle:       row.Str(metaTitle),
		MetaDescription: row.Str(metaDesc),
		Keywords:        row.Str(keywords),
		Primary:         row.Bool(isPrimary),
		Published:       row.Bool(published),
		Active:          row.Bool(active),
		Slug:            row.Str(slug),
		RequireAuth:     row.Bool(requireAuth),
		Canonical:       row.Str(canonical),
		Revisions:       <-revCh,
	}

	ch <- content
}

func GetPrimaryMenu() (menu Menu, err error) {
	sel, err := database.GetStatement("getPrimaryMenuStmt")
	if err != nil {
		return menu, err
	}
	row, res, err := sel.ExecFirst()
	if database.MysqlError(err) {
		return menu, err
	}
	id := row.Int(res.Map("menuID"))

	c := make(chan Menu)
	mi := make(chan []MenuItem)

	go PopulateMenu(row, res, c)
	go GetMenuItems(id, mi)
	menu = <-c
	menu.Items = <-mi

	return menu, nil
}

func PopulateMenu(row mysql.Row, res mysql.Result, ch chan Menu) {
	id := res.Map("menuID")
	name := res.Map("menu_name")
	isPrimary := res.Map("isPrimary")
	active := res.Map("active")
	displayName := res.Map("display_name")
	reqAuth := res.Map("requireAuthentication")
	showOnSitemap := res.Map("showOnSitemap")
	sort := res.Map("sort")
	var menu Menu

	menu = Menu{
		ID:            row.Int(id),
		Name:          row.Str(name),
		Primary:       row.Bool(isPrimary),
		Active:        row.Bool(active),
		DisplayName:   row.Str(displayName),
		RequireAuth:   row.Bool(reqAuth),
		ShowOnSitemap: row.Bool(showOnSitemap),
		Sort:          row.Int(sort),
	}

	ch <- menu
}

func GetMenuItems(id int, ch chan []MenuItem) {
	items := make([]MenuItem, 0)
	if id > 0 {
		sel, err := database.GetStatement("getMenuItemsStmt")
		if err != nil {
			ch <- items
			return
		}
		sel.Bind(id)
		rows, res, err := sel.Exec()
		if database.MysqlError(err) {
			ch <- items
			return
		}

		mch := make(chan MenuItem)
		for _, row := range rows {
			go PopulateMenuItem(row, res, mch)
		}
		for _, _ = range rows {
			items = append(items, <-mch)
		}

	}

	ch <- items
}

func PopulateMenuItem(row mysql.Row, res mysql.Result, ch chan MenuItem) {
	var item MenuItem

	mcid := res.Map("menuContentID")
	menuID := res.Map("menuID")
	contentID := res.Map("contentID")
	menuSort := res.Map("menuSort")
	menuTitle := res.Map("menuTitle")
	menuLink := res.Map("menuLink")
	parentID := res.Map("parentID")
	linkTarget := res.Map("linkTarget")

	id := row.Int(mcid)
	if id > 0 {
		cid := row.Int(contentID)
		cch := make(chan Content)
		go PopulateContent(row, res, cch)

		item = MenuItem{
			ID:         id,
			MenuID:     row.Int(menuID),
			ContentID:  cid,
			Sort:       row.Int(menuSort),
			Title:      row.Str(menuTitle),
			Link:       row.Str(menuLink),
			ParentID:   row.Int(parentID),
			LinkTarget: row.Str(linkTarget),
			Content:    <-cch,
		}
	}

	ch <- item
}

func GetContentRevisions(id int, ch chan []ContentRevision) {
	revisions := make([]ContentRevision, 0)
	if id > 0 {
		sel, err := database.GetStatement("GetContentRevisionsStmt")
		if err != nil {
			ch <- revisions
			return
		}
		sel.Bind(id)
		rows, res, err := sel.Exec()
		if database.MysqlError(err) {
			ch <- revisions
			return
		}

		rch := make(chan ContentRevision)
		for _, row := range rows {
			go PopulateRevision(row, res, rch)
		}
		for _, _ = range rows {
			revisions = append(revisions, <-rch)
		}

	}

	ch <- revisions
}

func PopulateRevision(row mysql.Row, res mysql.Result, ch chan ContentRevision) {
	var revision ContentRevision

	id := res.Map("revisionID")
	contentID := res.Map("contentID")
	contentText := res.Map("content_text")
	createdOn := res.Map("createdOn")
	active := res.Map("active")

	revision = ContentRevision{
		ID:          row.Int(id),
		ContentID:   row.Int(contentID),
		ContentText: row.Str(contentText),
		CreatedOn:   row.Time(createdOn, UTC),
		Active:      row.Bool(active),
	}

	ch <- revision
}
