package website

import (
	"../../helpers/plate"
	"../../models"
	"encoding/json"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var store = sessions.NewCookieStore([]byte("adminstuffs"))

func Index(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)

	contents, _ := models.GetAllSiteContent()
	menu, _ := models.GetPrimaryMenu()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "Mon, 01/02/06, 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}

	tmpl.Bag["PageTitle"] = "Main Menu"
	tmpl.Bag["menu"] = menu
	tmpl.Bag["contents"] = contents

	tmpl.HtmlTemplate.Parse(menu.GenerateHtml())

	tmpl.FuncMap["addPublishedClass"] = func(itm models.MenuItem) bool {
		if (itm.HasContent() && itm.Content.Published) || !itm.HasContent() {
			return true
		}
		return false
	}
	tmpl.FuncMap["equalsOne"] = func(num int) bool {
		return num == 1
	}

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/website/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Menus(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	menus, _ := models.GetAllMenus()
	tmpl.Bag["menus"] = menus
	tmpl.Bag["PageTitle"] = "Menus"

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/website/menus.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Menu(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)

	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	menu := models.Menu{ID: id}
	contents, _ := models.GetAllSiteContent()
	err = menu.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "Mon, 01/02/06, 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}

	tmpl.Bag["PageTitle"] = menu.Name + " Menu Contents"
	tmpl.Bag["menu"] = menu
	tmpl.Bag["contents"] = contents

	tmpl.HtmlTemplate.Parse(menu.GenerateHtml())

	tmpl.FuncMap["addPublishedClass"] = func(itm models.MenuItem) bool {
		if (itm.HasContent() && itm.Content.Published) || !itm.HasContent() {
			return true
		}
		return false
	}
	tmpl.FuncMap["equalsOne"] = func(num int) bool {
		return num == 1
	}

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/website/index.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	params := r.URL.Query()
	error := params.Get("error")
	menu := models.Menu{}

	tmpl.FuncMap["isZero"] = func(num int) bool {
		return num == 0
	}

	if len(error) > 0 {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Add Menu"
	tmpl.Bag["menu"] = menu

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/website/menuform.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	params := r.URL.Query()
	error := params.Get("error")
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	menu := models.Menu{ID: id}
	err = menu.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.FuncMap["isZero"] = func(num int) bool {
		return num == 0
	}

	if len(error) > 0 {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Edit " + menu.Name + " Menu Details"
	tmpl.Bag["menu"] = menu

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/website/menuform.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Save(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":id"))
	requireAuth, _ := strconv.ParseBool(r.FormValue("requireAuthentication"))
	showOnSitemap, _ := strconv.ParseBool(r.FormValue("showOnSitemap"))

	menu := models.Menu{
		ID:            id,
		Name:          r.FormValue("menu_name"),
		DisplayName:   r.FormValue("display_name"),
		RequireAuth:   requireAuth,
		ShowOnSitemap: showOnSitemap,
	}

	err := menu.Save()
	if err != nil {
		if menu.ID > 0 {
			http.Redirect(w, r, "/Website/Menu/Edit/"+strconv.Itoa(menu.ID)+"?error="+err.Error(), http.StatusFound)
		} else {
			http.Redirect(w, r, "/Website/Menu/Add?error="+err.Error(), http.StatusFound)
		}
	}
	http.Redirect(w, r, "/Website/Menus", http.StatusFound)
}

func Remove(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	success := false
	id, err := strconv.Atoi(params.Get(":id"))
	if err == nil {
		m := models.Menu{ID: id}
		err = m.Remove()
		if err == nil {
			success = true
		}
	}
	successobj := struct {
		Success bool
	}{
		success,
	}
	plate.ServeFormatted(w, r, successobj)
}

func SetPrimaryMenu(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err == nil {
		m := models.Menu{ID: id}
		err = m.SetPrimary()
	}
	http.Redirect(w, r, "/Website/Menus", http.StatusFound)
}

func MenuSort(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get("id"))
	if err == nil {
		menu := models.Menu{ID: id}
		pages := params["page"]
		menu.UpdateSort(pages)
	}
	plate.ServeFormatted(w, r, "")
}

func AddContentToMenu(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	menuID, err := strconv.Atoi(params.Get("menuid"))
	if err != nil {
		log.Println(err)
		return
	}
	contentID, err := strconv.Atoi(params.Get("contentid"))
	if err != nil {
		log.Println(err)
		return
	}
	menu := models.Menu{ID: menuID}
	menupage := menu.AddContent(contentID)
	menustuff := struct {
		MenuContentID int
		MenuSort      int
		ParentID      int
		ContentID     int
		PageTitle     string
		Published     bool
	}{
		MenuContentID: menupage.ID,
		MenuSort:      menupage.Sort,
		ParentID:      menupage.ParentID,
		ContentID:     menupage.ContentID,
		PageTitle:     menupage.Content.PageTitle,
		Published:     menupage.Content.Published,
	}
	plate.ServeFormatted(w, r, menustuff)

}

func RemoveContentAjax(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err == nil {
		models.RemoveContentFromMenu(id)
	}

	plate.ServeFormatted(w, r, "")
}

func SetPrimaryContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err == nil {
		menuid, err := strconv.Atoi(params.Get(":menuid"))
		if err == nil {
			models.SetPrimaryContent(id)
			http.Redirect(w, r, "/Website/Menu/"+strconv.Itoa(menuid), http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, "/Website", http.StatusFound)
}

func AddLink(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	error, _ := url.QueryUnescape(params.Get("error"))
	if len(strings.TrimSpace(error)) > 0 {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Add Link To Menu"
	tmpl.Bag["menuID"] = id
	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/website/addlink.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func EditLink(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	error, _ := url.QueryUnescape(params.Get("error"))
	if len(strings.TrimSpace(error)) > 0 {
		tmpl.Bag["error"] = error
	}
	mi := models.MenuItem{ID: id}
	item, err := mi.Get()
	if err != nil {
		tmpl.Bag["error"] = err.Error()
	}
	tmpl.Bag["PageTitle"] = "Edit Link"
	tmpl.Bag["item"] = item
	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/website/editlink.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func SaveNewLink(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	name := strings.TrimSpace(r.FormValue("link_name"))
	value := strings.TrimSpace(r.FormValue("link_value"))
	target, _ := strconv.ParseBool(r.FormValue("link_target"))
	if name == "" || value == "" {
		http.Redirect(w, r, "/Website/Link/Add/"+strconv.Itoa(id)+"?error="+url.QueryEscape("Title and Value are required"), http.StatusFound)
		return
	}
	item := models.MenuItem{
		MenuID:     id,
		Title:      name,
		Link:       value,
		LinkTarget: target,
	}
	err = item.SaveLink()
	if err != nil {
		http.Redirect(w, r, "/Website/Link/Add/"+strconv.Itoa(id)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
		return
	}
	http.Redirect(w, r, "/Website/Menu/"+strconv.Itoa(id), http.StatusFound)
}

func SaveLink(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	name := strings.TrimSpace(r.FormValue("link_name"))
	value := strings.TrimSpace(r.FormValue("link_value"))
	target, _ := strconv.ParseBool(r.FormValue("link_target"))
	if name == "" || value == "" {
		http.Redirect(w, r, "/Website/Link/Edit/"+strconv.Itoa(id)+"?error="+url.QueryEscape("Title and Value are required"), http.StatusFound)
		return
	}
	item := models.MenuItem{
		ID:         id,
		Title:      name,
		Link:       value,
		LinkTarget: target,
	}
	err = item.SaveLink()
	if err != nil {
		http.Redirect(w, r, "/Website/Link/Edit/"+strconv.Itoa(id)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
		return
	}
	item, err = item.Get()
	if err != nil {
		http.Redirect(w, r, "/Website/Link/Edit/"+strconv.Itoa(id)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
		return
	}
	http.Redirect(w, r, "/Website/Menu/"+strconv.Itoa(item.MenuID), http.StatusFound)

}

func CheckContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content := models.Content{ID: id}
	menujson := make([]struct{ MenuName string }, 0)
	menus := content.Check()
	if len(menus) > 0 {
		for _, menu := range menus {
			m := struct {
				MenuName string
			}{MenuName: menu}
			menujson = append(menujson, m)
		}
	}
	plate.ServeFormatted(w, r, menujson)
}

func DeleteContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content := models.Content{ID: id}
	successobj := struct{ Success bool }{Success: content.Delete()}

	plate.ServeFormatted(w, r, successobj)
}

func AddContent(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	error, _ := url.QueryUnescape(params.Get("error"))
	if len(strings.TrimSpace(error)) > 0 {
		tmpl.Bag["error"] = error
	}
	tmpl.FuncMap["isNotZero"] = func(num int) bool {
		return num != 0
	}
	content := models.Content{}
	session, _ := store.Get(r, "adminstuffs")
	if contents := session.Flashes("content"); len(contents) > 0 {
		json.Unmarshal([]byte(contents[0].(string)), &content)
	}
	pagecontent := ""
	if htmlcontent := session.Flashes("htmlcontent"); len(htmlcontent) > 0 {
		pagecontent = htmlcontent[0].(string)
	}
	session.Save(r, w)

	tmpl.Bag["PageTitle"] = "Add Content"
	tmpl.Bag["content"] = content
	tmpl.Bag["pagecontent"] = pagecontent
	tmpl.Bag["menuID"] = id
	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/website/addcontent.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}

}

func SaveNewContent(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("menuid"))
	reqauth, _ := strconv.ParseBool(r.FormValue("requireAuthentication"))
	publish, _ := strconv.ParseBool(r.FormValue("publish"))
	addtomenu, _ := strconv.ParseBool(r.FormValue("addtomenu"))
	pagecontent := r.FormValue("page_content")
	content := models.Content{
		PageTitle:       r.FormValue("page_title"),
		Keywords:        r.FormValue("keywords"),
		MetaTitle:       r.FormValue("meta_title"),
		MetaDescription: r.FormValue("meta_description"),
		Canonical:       r.FormValue("canonical"),
		RequireAuth:     reqauth,
		Published:       publish,
	}
	revision := models.ContentRevision{
		ContentText: pagecontent,
	}
	err := content.Save(revision)
	if err != nil {
		cjson, _ := json.Marshal(&content)
		session, _ := store.Get(r, "adminstuffs")
		session.AddFlash(string(cjson), "content")
		session.AddFlash(pagecontent, "htmlcontent")
		session.Save(r, w)
		http.Redirect(w, r, "/Website/Content/Add/"+strconv.Itoa(id)+"?error="+url.QueryEscape("Page Title is required"), http.StatusFound)
		return
	}
	if addtomenu && id > 0 {
		menu := models.Menu{ID: id}
		menu.AddContent(content.ID)
	}
	http.Redirect(w, r, "/Website/Menu/"+strconv.Itoa(id), http.StatusFound)
}

func EditContent(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	revid, err := strconv.Atoi(params.Get(":revid"))
	if err != nil {
		revid = 0
	}
	error, _ := url.QueryUnescape(params.Get("error"))
	if len(strings.TrimSpace(error)) > 0 {
		tmpl.Bag["error"] = error
	}
	message, _ := url.QueryUnescape(params.Get("message"))
	if len(strings.TrimSpace(message)) > 0 {
		tmpl.Bag["message"] = message
	}
	tmpl.FuncMap["isNotZero"] = func(num int) bool {
		return num != 0
	}

	content := models.Content{ID: id}
	content, err = content.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "Mon, 01/02/06, 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	revision := content.ActiveRevision
	if revid != 0 {
		revision = content.Revisions.GetRevision(revid)
	}
	tmpl.Bag["PageTitle"] = "Edit Content"
	tmpl.Bag["content"] = content
	tmpl.Bag["revision"] = revision
	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/website/editcontent.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}

}

func SaveContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":id"))
	revid, _ := strconv.Atoi(params.Get(":revid"))

	reqauth, _ := strconv.ParseBool(r.FormValue("requireAuthentication"))
	publish, _ := strconv.ParseBool(r.FormValue("publish"))
	pagecontent := r.FormValue("page_content")
	content := models.Content{
		ID:              id,
		PageTitle:       r.FormValue("page_title"),
		Keywords:        r.FormValue("keywords"),
		MetaTitle:       r.FormValue("meta_title"),
		MetaDescription: r.FormValue("meta_description"),
		Canonical:       r.FormValue("canonical"),
		RequireAuth:     reqauth,
		Published:       publish,
	}
	revision := models.ContentRevision{
		ID:          revid,
		ContentID:   id,
		ContentText: pagecontent,
	}
	err := content.Save(revision)
	if err != nil {
		cjson, _ := json.Marshal(&content)
		session, _ := store.Get(r, "adminstuffs")
		session.AddFlash(string(cjson), "content")
		session.AddFlash(pagecontent, "htmlcontent")
		session.Save(r, w)
		http.Redirect(w, r, "/Website/Content/Add/"+strconv.Itoa(id)+"?error="+url.QueryEscape("Page Title is required"), http.StatusFound)
		return
	}
	http.Redirect(w, r, "/Website/Content/Edit/"+strconv.Itoa(id)+"/"+strconv.Itoa(revid)+"?message="+url.QueryEscape("Content Page Updated Successfully!"), http.StatusFound)
}

func CopyRevision(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	revision := models.ContentRevision{ID: id}
	err = revision.Copy()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/Website/Content/Edit/"+strconv.Itoa(revision.ContentID), http.StatusFound)
}

func ActivateRevision(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	revision := models.ContentRevision{ID: id}
	err = revision.Activate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/Website/Content/Edit/"+strconv.Itoa(revision.ContentID), http.StatusFound)
}

func DeleteRevision(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	revision := models.ContentRevision{ID: id}
	err = revision.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/Website/Content/Edit/"+strconv.Itoa(revision.ContentID), http.StatusFound)
}
