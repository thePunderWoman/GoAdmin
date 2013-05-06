package faq

import (
	"../../helpers/plate"
	"../../models"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	faqs, _ := models.FAQ{}.GetAll()

	tmpl.Bag["PageTitle"] = "FAQ"
	tmpl.Bag["faqs"] = faqs

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/faq/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	error := r.URL.Query().Get("error")
	answer, _ := url.QueryUnescape(r.URL.Query().Get("answer"))
	question, _ := url.QueryUnescape(r.URL.Query().Get("question"))

	faq := models.FAQ{Question: question, Answer: answer}

	if strings.TrimSpace(error) == "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Add a FAQ Question"
	tmpl.Bag["faq"] = faq

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/faq/form.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	error := r.URL.Query().Get("error")

	faq, _ := models.FAQ{ID: id}.Get()

	if strings.TrimSpace(error) == "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Edit a FAQ Question"
	tmpl.Bag["faq"] = faq

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/faq/form.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Save(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	faq := models.FAQ{
		ID:       id,
		Question: r.FormValue("question"),
		Answer:   r.FormValue("answer"),
	}
	err := faq.Save()
	if err != nil {
		if id > 0 {
			http.Redirect(w, r, "/FAQ/Edit/"+strconv.Itoa(id)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
		} else {
			http.Redirect(w, r, "/FAQ/Add?error="+url.QueryEscape(err.Error())+"&question="+url.QueryEscape(faq.Question)+"&answer="+url.QueryEscape(faq.Answer), http.StatusFound)
		}
		return
	}
	http.Redirect(w, r, "/FAQ", http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))

	faq := models.FAQ{ID: id}
	err := faq.Delete()

	if err == nil {
		plate.ServeFormatted(w, r, "")
	} else {
		plate.ServeFormatted(w, r, err.Error())
	}
}
