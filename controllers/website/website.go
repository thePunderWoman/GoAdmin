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
	tmpl.Bag["displaymenu"] = menu.GenerateDisplayStructure()
	tmpl.Bag["contents"] = contents

	tmpl.FuncMap["addPublishedClass"] = func(itm models.MenuItem) bool {
		if (itm.HasContent() && itm.Content.Published) || !itm.HasContent() {
			return true
		}
		return false
	}
	tmpl.FuncMap["incrementCounter"] = func(num int) int {
		return num + 1
	}
	tmpl.FuncMap["equalsOne"] = func(num int) bool {
		return num == 1
	}

	tmpl.ParseFile("templates/website/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}
