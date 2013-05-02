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
	tmpl.Bag["menuID"] = id
	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/website/addlink.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
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
		http.Redirect(w, r, "/Website/Link/Add/"+strconv.Itoa(id)+"?error="+url.QueryEscape("Title and Value are required"), http.StatusFound)
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
	}
	http.Redirect(w, r, "/Website/Menu/"+strconv.Itoa(id), http.StatusFound)
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

func SaveContent(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("menuid"))
	reqauth, _ := strconv.ParseBool(r.FormValue("requireAuthentication"))
	content := models.Content{
		PageTitle:       r.FormValue("page_title"),
		Keywords:        r.FormValue("keywords"),
		MetaTitle:       r.FormValue("meta_title"),
		MetaDescription: r.FormValue("meta_description"),
		Canonical:       r.FormValue("canonical"),
		RequireAuth:     reqauth,
	}
	cjson, _ := json.Marshal(&content)
	session, _ := store.Get(r, "adminstuffs")
	session.AddFlash(string(cjson), "content")
	session.AddFlash(r.FormValue("page_content"), "htmlcontent")
	session.Save(r, w)
	http.Redirect(w, r, "/Website/Content/Add/"+strconv.Itoa(id), http.StatusFound)
}
