package customer

import (
	"../../helpers/plate"
	"../../models"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	//"net/url"
	//"strconv"
	//"strings"
	"time"
)

var store = sessions.NewCookieStore([]byte("adminstuffs"))

func Index(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)

	customers, _ := models.Customer{}.GetAll()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "Mon, 01/02/06, 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.Bag["PageTitle"] = "View Customers"
	tmpl.Bag["customers"] = customers

	tmpl.ParseFile("templates/customer/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}
