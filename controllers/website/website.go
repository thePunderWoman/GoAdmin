package website

import (
	"../../helpers/plate"
	"../../models"
	"log"
	"net/http"
	_ "net/url"
	_ "strconv"
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

	tmpl.ParseFile("templates/website/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}
