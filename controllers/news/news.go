package news

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

	news, _ := models.NewsItem{}.GetAll()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "Mon, 01/02/06, 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.Bag["PageTitle"] = "News"
	tmpl.Bag["news"] = news

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/news/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	error := r.URL.Query().Get("error")

	newsitem := models.NewsItem{}

	if strings.TrimSpace(error) == "" {
		tmpl.Bag["error"] = error
	}
	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 3:04 pm"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.Bag["PageTitle"] = "Add News Item"
	tmpl.Bag["newsitem"] = newsitem

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/news/form.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	error := r.URL.Query().Get("error")

	newsitem, _ := models.NewsItem{ID: id}.Get()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 3:04 pm"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	if strings.TrimSpace(error) == "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Edit News Item"
	tmpl.Bag["newsitem"] = newsitem

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/news/form.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Save(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	loc, _ := time.LoadLocation("UTC")
	pubstart, _ := time.Parse("01/02/2006 3:04 pm", r.FormValue("publishStart"))
	pubend, _ := time.Parse("01/02/2006 3:04 pm", r.FormValue("publishEnd"))
	newsitem := models.NewsItem{
		ID:           id,
		Title:        r.FormValue("title"),
		Lead:         r.FormValue("lead"),
		Content:      r.FormValue("content"),
		PublishStart: ChangeZone(pubstart, loc),
		PublishEnd:   ChangeZone(pubend, loc),
	}
	err := newsitem.Save()
	if err != nil {
		if id > 0 {
			http.Redirect(w, r, "/News/Edit/"+strconv.Itoa(id)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
		} else {
			http.Redirect(w, r, "/News/Add?error="+url.QueryEscape(err.Error()), http.StatusFound)
		}
		return
	}
	http.Redirect(w, r, "/News", http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))

	newsitem := models.NewsItem{ID: id}
	err := newsitem.Delete()

	if err == nil {
		plate.ServeFormatted(w, r, "")
	} else {
		plate.ServeFormatted(w, r, err.Error())
	}
}

func ChangeZone(t time.Time, zone *time.Location) time.Time {
	if !t.IsZero() {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), zone)
	}
	return t
}
