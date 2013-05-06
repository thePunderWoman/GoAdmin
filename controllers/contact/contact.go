package contact

import (
	"../../helpers/plate"
	"../../models"
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

func EditReceiver(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":id"))
	error, _ := url.QueryUnescape(params.Get("error"))

	receiver, err := models.ContactReceiver{ID: id}.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	types, err := models.ContactType{}.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.FuncMap["hasType"] = func(typeID int) bool {
		for _, t := range receiver.Types {
			if typeID == t.ID {
				return true
			}
		}
		return false
	}

	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Edit Contact Receiver"
	tmpl.Bag["types"] = types
	tmpl.Bag["receiver"] = receiver

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/contact/receiver.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func AddReceiver(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	params := r.URL.Query()
	error, _ := url.QueryUnescape(params.Get("error"))

	receiver := models.ContactReceiver{}
	types, err := models.ContactType{}.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.FuncMap["hasType"] = func(typeID int) bool {
		for _, t := range receiver.Types {
			if typeID == t.ID {
				return true
			}
		}
		return false
	}

	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Add Contact Receiver"
	tmpl.Bag["types"] = types
	tmpl.Bag["receiver"] = receiver

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/contact/receiver.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func SaveReceiver(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("receiverID"))
	if err != nil {
		id = 0
	}
	receiver := models.ContactReceiver{
		ID:        id,
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
	}
	types := r.Form["types"]
	err = receiver.Save(types)
	if err != nil {
		if id == 0 {
			http.Redirect(w, r, "/Contact/Receivers", http.StatusFound)
		} else {
			http.Redirect(w, r, "/Contact/EditReceiver"+strconv.Itoa(id)+url.QueryEscape(err.Error()), http.StatusFound)

		}
	}
	http.Redirect(w, r, "/Contact/Receivers", http.StatusFound)
}

func DeleteReceiver(w http.ResponseWriter, r *http.Request) {
	urlsuffix := ""
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	receiver := models.ContactReceiver{ID: id}

	err = receiver.Delete()
	if err != nil {
		urlsuffix = "?error=" + url.QueryEscape(err.Error())
	}
	http.Redirect(w, r, "/Contact/Receivers"+urlsuffix, http.StatusFound)
}

func AddType(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	params := r.URL.Query()
	error, _ := url.QueryUnescape(params.Get("error"))
	name, _ := url.QueryUnescape(params.Get("name"))

	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["name"] = name
	tmpl.Bag["PageTitle"] = "Add Contact Type"

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/contact/addtype.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func SaveType(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Redirect(w, r, "/Contact/AddType?error="+url.QueryEscape("name is required"), http.StatusFound)
		return
	}
	typeobj := models.ContactType{Name: name}
	err := typeobj.Save()
	if err != nil {
		http.Redirect(w, r, "/Contact/AddType?name="+url.QueryEscape(name)+"&error="+url.QueryEscape(err.Error()), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/Contact/Types", http.StatusFound)
}

func DeleteType(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	typeobj := models.ContactType{ID: id}
	err := typeobj.Delete()
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/Contact/Types", http.StatusFound)
}
