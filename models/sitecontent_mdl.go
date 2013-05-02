package models

import (
	"../helpers/database"
	_ "errors"
	"github.com/ziutek/mymysql/mysql"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
	LinkTarget bool
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
	Revisions       ContentRevisions
	ActiveRevision  ContentRevision
}
type ContentRevision struct {
	ID          int
	ContentID   int
	ContentText string
	CreatedOn   time.Time
	Active      bool
}
type ContentRevisions []ContentRevision

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
	revCh := make(chan ContentRevisions)
	go GetContentRevisions(id, revCh)
	revisions := <-revCh
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
		Revisions:       revisions,
		ActiveRevision:  revisions.GetActiveRevision(),
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
			LinkTarget: row.Bool(linkTarget),
			Content:    <-cch,
		}
	}

	ch <- item
}

func (i *MenuItem) HasContent() bool {
	return i.ContentID > 0
}

func GetContentRevisions(id int, ch chan ContentRevisions) {
	var revisions ContentRevisions
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
	revisions.Sort()
	ch <- revisions
}

func (r ContentRevisions) GetActiveRevision() ContentRevision {
	for _, revision := range r {
		if revision.Active {
			return revision
		}
	}
	var active ContentRevision
	return active
}

func (r ContentRevisions) GetRevision(id int) ContentRevision {
	for _, revision := range r {
		if revision.ID == id {
			return revision
		}
	}
	var active ContentRevision
	return active
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

func (m *Menu) UpdateSort(pages []string) {
	for _, page := range pages {
		pdata := strings.Split(page, "-")
		upd, err := database.GetStatement("updateMenuItemSortStmt")
		if err != nil {
			return
		}
		sort, _ := strconv.Atoi(pdata[2])
		parentID, _ := strconv.Atoi(pdata[1])
		id, _ := strconv.Atoi(pdata[0])

		params := struct {
			Sort     int
			ParentID int
			ID       int
		}{
			Sort:     sort,
			ParentID: parentID,
			ID:       id,
		}

		upd.Bind(&params)

		_, _, err = upd.Exec()
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (m *Menu) AddContent(contentID int) (item MenuItem) {
	sort, err := m.GetNextSort()
	if err != nil {
		return item
	}

	// run insert statement
	ins, err := database.GetStatement("addMenuContentItemStmt")
	if err != nil {
		return item
	}
	ins.Bind(m.ID, contentID, sort)
	_, res, err := ins.Exec()
	if err != nil {
		log.Println(err)
		return item
	}
	id := res.InsertId()

	if id > 0 {
		// get item
		sel3, err := database.GetStatement("getMenuItemStmt")
		if err != nil {
			return item
		}
		sel3.Bind(id)
		row, res, err := sel3.ExecFirst()
		if err != nil {
			log.Println(err)
			return item
		}
		mch := make(chan MenuItem)
		go PopulateMenuItem(row, res, mch)
		item = <-mch
	}

	return item
}

func RemoveContentFromMenu(id int) {
	sel1, err := database.GetStatement("getMenuItemStmt")
	if err != nil {
		log.Println(err)
		return
	}
	sel1.Bind(id)
	row, res, err := sel1.ExecFirst()
	if err != nil {
		log.Println(err)
		return
	}
	mch := make(chan MenuItem)
	go PopulateMenuItem(row, res, mch)
	item := <-mch

	sel2, err := database.GetStatement("GetMenuParentsStmt")
	if err != nil {
		log.Println(err)
		return
	}
	sel2.Reset()
	sel2.Bind(item.ParentID, item.MenuID)
	rows, res, err := sel2.Exec()
	if err != nil {
		log.Println(err)
		return
	}
	mID := res.Map("menuContentID")
	pID := res.Map("parentID")
	pitems := make([]MenuItem, 0)
	citems := make([]MenuItem, 0)

	for _, row = range rows {
		mitem := MenuItem{
			ID:       row.Int(mID),
			ParentID: row.Int(pID),
		}
		pitems = append(pitems, mitem)
	}

	sel3, err := database.GetStatement("GetMenuParentsStmt")
	if err != nil {
		log.Println(err)
		return
	}
	sel3.Reset()
	sel3.Bind(item.ID, item.MenuID)
	rows, res, err = sel3.Exec()
	if err != nil {
		log.Println(err)
		return
	}

	mID = res.Map("menuContentID")
	pID = res.Map("parentID")

	for _, row = range rows {
		mitem := MenuItem{
			ID:       row.Int(mID),
			ParentID: row.Int(pID),
		}
		citems = append(citems, mitem)
	}

	sort := 0
	for _, itm := range pitems {
		sort += 1
		if item.ID == itm.ID {
			// adjust children
			for _, citm := range citems {
				sort += 1
				citm.ParentID = itm.ParentID
				citm.Sort = sort
				go citm.SetSort()
			}
		} else {
			// adjust sort
			itm.Sort = sort
			go itm.SetSort()
		}
	}

	del, err := database.GetStatement("DeleteMenuItemStmt")
	if err != nil {
		log.Println(err)
		return
	}
	del.Reset()
	del.Bind(item.ID)
	_, _, err = del.Exec()
	if err != nil {
		log.Println(err)
	}
	return
}

func (m *Menu) GetNextSort() (int, error) {
	sel, err := database.GetStatement("getMenuSortStmt")
	if err != nil {
		return 0, err
	}
	sel.Bind(m.ID)
	row, res, err := sel.ExecFirst()
	if database.MysqlError(err) {
		return 0, err
	}
	msort := res.Map("menuSort")
	sort := row.Int(msort)
	sort += 1
	return sort, nil
}

func (m *MenuItem) SetSort() {
	upd, err := database.GetStatement("updateMenuItemSortStmt")
	if err != nil {
		return
	}
	upd.Bind(m.Sort, m.ParentID, m.ID)
	_, _, err = upd.Exec()
	if err != nil {
		log.Println(err)
	}
	return
}

func (m *MenuItem) SaveLink() error {
	menu := Menu{ID: m.MenuID}
	sort, err := menu.GetNextSort()
	if err != nil {
		return err
	}

	ins, err := database.GetStatement("addMenuLinkItemStmt")
	if err != nil {
		return err
	}

	ins.Bind(m.MenuID, m.Title, m.Link, m.LinkTarget, sort)
	_, _, err = ins.Exec()
	if err != nil {
		return err
	}
	return nil
}

func SetPrimaryContent(id int) {
	upd, err := database.GetStatement("clearPrimaryContentStmt")
	if err != nil {
		return
	}
	_, _, err = upd.Exec()
	if err != nil {
		log.Println(err)
	}
	set, err := database.GetStatement("setPrimaryContentStmt")
	if err != nil {
		return
	}
	set.Bind(id)
	_, _, err = set.Exec()
	if err != nil {
		log.Println(err)
	}
	return
}

func (m *MenuItems) SortItems() {
	sort.Sort(m)
}

func (r *ContentRevisions) Sort() {
	sort.Sort(r)
}

func (c *Content) Check() (names []string) {
	sel, err := database.GetStatement("checkContentStmt")
	if err != nil {
		log.Println(err)
		return names
	}
	sel.Bind(c.ID)
	rows, res, err := sel.Exec()
	if err != nil {
		log.Println(err)
		return names
	}
	name := res.Map("menu_name")
	for _, row := range rows {
		names = append(names, row.Str(name))
	}

	return names
}

func (c *Content) Get() (Content, error) {
	var content Content
	sel, err := database.GetStatement("getContentStmt")
	if err != nil {
		return content, err
	}
	sel.Bind(c.ID)
	row, res, err := sel.ExecFirst()
	if err != nil {
		return content, err
	}
	ch := make(chan Content)
	go PopulateContent(row, res, ch)
	content = <-ch
	return content, nil
}

func (c *Content) Delete() bool {
	del, err := database.GetStatement("deleteContentStmt")
	if err != nil {
		log.Println(err)
		return false
	}
	del.Bind(c.ID)
	_, _, err = del.Exec()
	if err != nil {
		return false
	}

	return true
}

func (c *Content) Save(pagecontent string) error {
	ins, err := database.GetStatement("addContentStmt")
	if err != nil {
		log.Println(err)
		return err
	}
	params := struct {
		Title       string
		Created     time.Time
		Modified    time.Time
		MetaTitle   string
		MetaDesc    string
		Keywords    string
		Published   bool
		Slug        string
		RequireAuth bool
		Canonical   string
	}{
		Title:       c.PageTitle,
		Created:     time.Now().In(UTC),
		Modified:    time.Now().In(UTC),
		MetaTitle:   c.MetaTitle,
		MetaDesc:    c.MetaDescription,
		Keywords:    c.Keywords,
		Published:   c.Published,
		Slug:        GenerateSlug(c.PageTitle),
		RequireAuth: c.RequireAuth,
		Canonical:   c.Canonical,
	}

	ins.Bind(&params)
	_, res, err := ins.Exec()
	if err != nil {
		return err
	}
	id := res.InsertId()
	c.ID = int(id)

	ins2, err := database.GetStatement("addContentRevisionStmt")
	if err != nil {
		log.Println(err)
		return err
	}
	ins2.Reset()
	ins2.Bind(c.ID, pagecontent, time.Now().In(UTC), true)
	_, _, err = ins2.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (m MenuItems) Len() int                  { return len(m) }
func (m MenuItems) Swap(i, j int)             { m[i], m[j] = m[j], m[i] }
func (m MenuItems) Less(i, j int) bool        { return m[i].Sort < m[j].Sort }
func (r ContentRevisions) Len() int           { return len(r) }
func (r ContentRevisions) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ContentRevisions) Less(i, j int) bool { return r[i].CreatedOn.Before(r[j].CreatedOn) }

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	invchars := regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
	underscores := regexp.MustCompile(`\s+`)

	slug = invchars.ReplaceAllString(slug, "")
	slug = underscores.ReplaceAllString(slug, "_")

	return slug
}
