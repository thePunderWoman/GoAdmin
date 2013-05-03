package contact

import (
	"../../helpers/plate"
	"../../models"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strconv"
	"time"
)

var store = sessions.NewCookieStore([]byte("adminstuffs"))

func Index(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)

	contacts, _ := models.Contact{}.GetAll()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "Mon, 01/02/06, 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.Bag["PageTitle"] = "View Contacts"
	tmpl.Bag["contacts"] = contacts

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/contact/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func View(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":id"))

	contact := models.Contact{ID: id}
	err := contact.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "Mon, 01/02/06, 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}

	tmpl.Bag["PageTitle"] = "View Contact Details"
	tmpl.Bag["contact"] = contact

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/contact/view.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Receivers(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	receivers, err := models.ContactReceiver{}.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Bag["PageTitle"] = "Contact Receivers"
	tmpl.Bag["receivers"] = receivers

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/contact/receivers.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Types(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	types, err := models.ContactType{}.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Bag["PageTitle"] = "Contact Types"
	tmpl.Bag["types"] = types

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/contact/types.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}
