package website

import (
	"../../helpers/plate"
	"../../models"
	"log"
	"net/http"
	_ "net/url"
	"strconv"
	_ "strings"
	"time"
)

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
	log.Println(menu.ID)
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
