package models

import (
	"../helpers/database"
	_ "errors"
	"github.com/ziutek/mymysql/mysql"
	_ "log"
	"sort"
	"strconv"
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
	Items         MenuItems
	Map           MenuMap
}

type MenuItems []MenuItem
type MenuMap struct {
	Items map[int]MenuItems
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

func (m *Menu) Get() error {
	sel, err := database.GetStatement("getMenuByIDStmt")
	if err != nil {
		return err
	}
	sel.Bind(m.ID)

	row, res, err := sel.ExecFirst()
	if database.MysqlError(err) {
		return err
	}

	c := make(chan Menu)
	mi := make(chan []MenuItem)

	go PopulateMenu(row, res, c)
	go GetMenuItems(m.ID, mi)
	menu := <-c

	m.ID = menu.ID
	m.Name = menu.Name
	m.DisplayName = menu.DisplayName
	m.Primary = menu.Primary
	m.RequireAuth = menu.RequireAuth
	m.Active = menu.Active
	m.ShowOnSitemap = menu.ShowOnSitemap
	m.Sort = menu.Sort
	m.Items = <-mi

	return nil
}

func GetAllMenus() (menus []Menu, err error) {
	sel, err := database.GetStatement("GetAllMenusStmt")
	if err != nil {
		return menus, err
	}
	rows, res, err := sel.Exec()
	if database.MysqlError(err) {
		return menus, err
	}

	c := make(chan Menu)

	for _, row := range rows {
		go PopulateMenu(row, res, c)
	}
	for _, _ = range rows {
		menus = append(menus, <-c)
	}

	return menus, nil

}

func (m *Menu) GenerateMenuMap() {
	var menumap MenuMap
	menumap.Items = make(map[int]MenuItems, 0)

	for _, item := range m.Items {
		if menumap.Items[item.ParentID] == nil {
			var mitems MenuItems
			menumap.Items[item.ParentID] = mitems
		}
		menumap.Items[item.ParentID] = append(menumap.Items[item.ParentID], item)
	}
	for _, tier := range menumap.Items {
		tier.SortItems()
	}
	m.Map = menumap
}

func (m *Menu) GenerateHtml() string {
	html := ""
	rootkids := ""
	m.GenerateMenuMap()
	html += `{{ define "menucontent" }}<ul id="pages" class="connected">`
	if _, ok := m.Map.Items[0]; ok {
		html += m.Map.Build(0, m.ID, 1)
		rootkids, _ = m.Map.GetChildren(0)
	}
	html += `</ul><input type="hidden" id="children_0" value="` + rootkids + `" />{{ end }}`
	return html
}

func (mm *MenuMap) GetChildren(parentID int) (string, int) {
	list := ""
	counter := 0
	for _, item := range mm.Items[parentID] {
		counter += 1
		if counter > 1 {
			list += ","
		}
		list += strconv.Itoa(item.ID)
	}
	return list, len(mm.Items[parentID])
}

func (mm *MenuMap) HasChildren(parentID int) bool {
	if _, ok := mm.Items[parentID]; ok {
		return true
	}
	return false
}

func (mm *MenuMap) Build(parentID int, menuID int, level int) string {
	html := ""
	for _, item := range mm.Items[parentID] {
		// generate menu

		html += `<li class="level_` + strconv.Itoa(level)
		if (item.HasContent() && item.Content.Published) || !item.HasContent() {
			html += " published"
		}
		html += `" id="item_` + strconv.Itoa(item.ID) + `"><span class="handle">↕</span>`
		html += `<span class="title">`
		if item.HasContent() {
			html += item.Content.PageTitle
		} else {
			html += item.Title + " (link)"
		}
		html += `</span><span class="controls">`
		if item.HasContent() && item.Content.Primary {
			html += `<a href="/Website/SetPrimaryContent/` + strconv.Itoa(item.ContentID) + `/` + strconv.Itoa(menuID) + `"><img src="/img/check.png" alt="Primary Page" title="Primary Page" /></a>`
		} else if item.HasContent() {
			html += `<a href="/Website/SetPrimaryContent/` + strconv.Itoa(item.ContentID) + `/` + strconv.Itoa(menuID) + `"><img src="/img/makeprimary.png" alt="Make This Page the Primary Page" title="Make This Page the Primary Page" /></a>`
		}
		if item.HasContent() {
			html += `<a href="/Website/Content/Edit/` + strconv.Itoa(item.ContentID) + `"><img src="/img/pencil.png" alt="Edit Page" title="Edit Page" /></a>`
		} else {
			html += `<a href="/Website/Link/Edit/` + strconv.Itoa(item.ID) + `"><img src="/img/pencil.png" alt="Edit Link" title="Edit Link" /></a>`
		}
		html += `<a href="/Website/RemoveContent/` + strconv.Itoa(item.ID) + `" class="remove" id="remove_` + strconv.Itoa(item.ID) + `"><img src="/img/delete.png" alt="Remove Page From Menu" title="Remove Page From Menu" /></a>`
		html += `</span><span id="meta_` + strconv.Itoa(item.ID) + `">`
		html += `<input type="hidden" id="parent_` + strconv.Itoa(item.ID) + `" value="` + strconv.Itoa(item.ParentID) + `" />`
		children, childcount := mm.GetChildren(item.ID)
		html += `<input type="hidden" id="children_` + strconv.Itoa(item.ID) + `" value="` + children + `" />`
		html += `<input type="hidden" id="count_` + strconv.Itoa(item.ID) + `" value="` + strconv.Itoa(childcount) + `" />`
		html += `<input type="hidden" id="sort_` + strconv.Itoa(item.ID) + `" value="` + strconv.Itoa(item.Sort) + `" />`
		html += `<input type="hidden" id="depth_` + strconv.Itoa(item.ID) + `" value="1" />`
		html += `</span><ul id="transport_` + strconv.Itoa(item.ID) + `"></ul></li>`
		if mm.HasChildren(item.ID) {
			html += mm.Build(item.ID, menuID, level+1)
		}

	}
	return html
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
	var items MenuItems
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

func (i *MenuItem) HasContent() bool {
	return i.ContentID > 0
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

func (m *Menu) Save() error {
	if m.ID > 0 {
		// update
		// new
		upd, err := database.GetStatement("UpdateMenuStmt")
		if err != nil {
			return err
		}

		params := struct {
			Name          string
			RequireAuth   bool
			ShowOnSitemap bool
			DisplayName   string
			ID            int
		}{}

		params.Name = m.Name
		params.DisplayName = m.DisplayName
		params.RequireAuth = m.RequireAuth
		params.ShowOnSitemap = m.ShowOnSitemap
		params.ID = m.ID

		upd.Bind(&params)

		_, _, err = upd.Exec()
		if err != nil {
			return err
		}

	} else {
		// new
		ins, err := database.GetStatement("AddMenuStmt")
		if err != nil {
			return err
		}

		params := struct {
			Name          string
			DisplayName   string
			RequireAuth   bool
			ShowOnSitemap bool
		}{}

		params.Name = m.Name
		params.DisplayName = m.DisplayName
		params.RequireAuth = m.RequireAuth
		params.ShowOnSitemap = m.ShowOnSitemap

		ins.Bind(&params)

		_, _, err = ins.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Menu) Remove() error {
	if m.ID > 0 {
		// update
		// new
		upd, err := database.GetStatement("deleteMenuStmt")
		if err != nil {
			return err
		}

		upd.Bind(m.ID)

		_, _, err = upd.Exec()
		if err != nil {
			return err
		}

	}
	return nil
}

func (m *Menu) SetPrimary() error {
	m.Get()
	upd, err := database.GetStatement("clearPrimaryMenuStmt")
	if err != nil {
		return err
	}

	_, _, err = upd.Exec()
	if err != nil {
		return err
	}

	if !m.Primary {
		upd1, err := database.GetStatement("setPrimaryMenuStmt")
		if err != nil {
			return err
		}

		upd1.Bind(m.ID)

		_, _, err = upd.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MenuItems) SortItems() {
	sort.Sort(m)
}

func (m MenuItems) Len() int           { return len(m) }
func (m MenuItems) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MenuItems) Less(i, j int) bool { return m[i].Sort < m[j].Sort }
