package users

import (
	"../../helpers/plate"
	"../../models"
	"log"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {

	server := plate.NewServer()

	tmpl, err := plate.GetTemplate()
	if err != nil {
		tmpl, err = server.Template(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	u := models.User{}
	users, err := u.GetAll()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "Mon, 01/02/06, 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}

	tmpl.Bag["Users"] = users
	tmpl.Bag["Count"] = len(users)

	tmpl.ParseFile("templates/users/index.html", false)
	log.Println(tmpl.HtmlTemplate)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}
