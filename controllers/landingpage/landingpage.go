package landingpage

import (
	"../../helpers/plate"
	"../../models"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	pages, _ := models.LandingPage{}.GetActive()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 3:04 PM"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.Bag["PageTitle"] = "Landing Pages"
	tmpl.Bag["pages"] = pages
	tmpl.Bag["type"] = 1
	tmpl.Bag["website"] = "http://www.curtmfg.com/"

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/landingpage/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Past(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	pages, _ := models.LandingPage{}.GetPast()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 3:04 PM"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.Bag["PageTitle"] = "Landing Pages"
	tmpl.Bag["pages"] = pages
	tmpl.Bag["type"] = 0
	tmpl.Bag["website"] = "http://www.curtmfg.com/"

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/landingpage/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	error, _ := url.QueryUnescape(r.URL.Query().Get("error"))

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 03:04 pm"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}

	tmpl.Bag["PageTitle"] = "Add a Landing Page"

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/landingpage/add.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func New(w http.ResponseWriter, r *http.Request) {
	start, _ := time.ParseInLocation("01/02/2006 03:04 pm", r.FormValue("startDate"), models.CentralTime)
	end, _ := time.ParseInLocation("01/02/2006 03:04 pm", r.FormValue("endDate"), models.CentralTime)
	newwindow, _ := strconv.ParseBool(r.FormValue("newWindow"))
	page := models.LandingPage{
		Name:         r.FormValue("name"),
		Start:        start,
		End:          end,
		URL:          r.FormValue("url"),
		NewWindow:    newwindow,
		MenuPosition: r.FormValue("menuPosition"),
	}
	err := page.New()
	if err != nil {
		http.Redirect(w, r, "/LandingPages/Add?error="+url.QueryEscape(err.Error()), http.StatusFound)
		return
	}
	http.Redirect(w, r, "/LandingPages/Edit/"+strconv.Itoa(page.ID), http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	error, _ := url.QueryUnescape(r.URL.Query().Get("error"))

	page, err := models.LandingPage{ID: id}.Get()
	if err != nil {
		log.Println(err)
	}

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 03:04 pm"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.FuncMap["equalsPosition"] = func(pos string) bool {
		return pos == page.MenuPosition
	}
	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Edit a Landing Page"
	tmpl.Bag["page"] = page

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/landingpage/edit.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Save(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	start, _ := time.ParseInLocation("01/02/2006 03:04 pm", r.FormValue("startDate"), models.CentralTime)
	end, _ := time.ParseInLocation("01/02/2006 03:04 pm", r.FormValue("endDate"), models.CentralTime)
	newwindow, _ := strconv.ParseBool(r.FormValue("newWindow"))
	hash := r.URL.Fragment
	page := models.LandingPage{
		ID:              id,
		Name:            r.FormValue("name"),
		Start:           start,
		End:             end,
		URL:             r.FormValue("url"),
		NewWindow:       newwindow,
		MenuPosition:    r.FormValue("menuPosition"),
		Content:         r.FormValue("page_content"),
		LinkClasses:     r.FormValue("linkClasses"),
		ConversionID:    r.FormValue("conversionID"),
		ConversionLabel: r.FormValue("conversionLabel"),
	}
	urlstr := "/LandingPages/Edit/" + strconv.Itoa(page.ID)
	err := page.Save()
	if err != nil {
		urlstr += "?error=" + url.QueryEscape(err.Error())
	}
	if strings.TrimSpace(hash) != "" {
		urlstr += "#" + hash
	}
	http.Redirect(w, r, urlstr, http.StatusFound)
}

func AddData(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("pageID"))
	data := models.LandingPageData{
		LandingPageID: id,
		Key:           r.URL.Query().Get("key"),
		Value:         r.URL.Query().Get("value"),
	}
	datalist := data.Save()
	plate.ServeFormatted(w, r, datalist)
}

func RemoveData(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))

	data := models.LandingPageData{ID: id}
	success := data.Delete()

	plate.ServeFormatted(w, r, success)
}

func Remove(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))

	page := models.LandingPage{ID: id}
	success := page.Delete()

	plate.ServeFormatted(w, r, success)
}
