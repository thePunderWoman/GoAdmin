package salesrep

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

	reps, _ := models.SalesRep{}.GetAll()

	tmpl.Bag["PageTitle"] = "Sales Reps"
	tmpl.Bag["reps"] = reps

	tmpl.ParseFile("templates/salesrep/navigation.html", false)
	tmpl.ParseFile("templates/salesrep/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	error := r.URL.Query().Get("error")

	var rep models.SalesRep

	tmpl.Bag["PageTitle"] = "Add Sales Rep"
	tmpl.Bag["Type"] = "Add"
	tmpl.Bag["rep"] = rep
	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}

	tmpl.ParseFile("templates/salesrep/navigation.html", false)
	tmpl.ParseFile("templates/salesrep/form.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	error := r.URL.Query().Get("error")

	rep, _ := models.SalesRep{ID: id}.Get()

	tmpl.Bag["PageTitle"] = "Edit Sales Rep"
	tmpl.Bag["Type"] = "Edit"
	tmpl.Bag["rep"] = rep
	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}

	tmpl.ParseFile("templates/salesrep/navigation.html", false)
	tmpl.ParseFile("templates/salesrep/form.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Save(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	rep := models.SalesRep{
		ID:   id,
		Name: r.FormValue("name"),
		Code: r.FormValue("code"),
	}
	err := rep.Save()
	if err != nil {
		if rep.ID > 0 {
			http.Redirect(w, r, "/SalesRep/Edit/"+strconv.Itoa(rep.ID)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
		} else {
			http.Redirect(w, r, "/SalesRep/Add?error="+url.QueryEscape(err.Error()), http.StatusFound)
		}
	}
	http.Redirect(w, r, "/SalesRep", http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	rep := models.SalesRep{ID: id}
	success := rep.Delete()
	plate.ServeFormatted(w, r, success)
}
