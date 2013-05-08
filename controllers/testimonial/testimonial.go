package testimonial

import (
	"../../helpers/plate"
	"../../models"
	"log"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	testimonials, _ := models.Testimonial{}.GetUnapproved()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 3:04 PM"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.Bag["PageTitle"] = "All UnApproved Testimonials"
	tmpl.Bag["testimonials"] = testimonials
	tmpl.Bag["type"] = "Unapproved"

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/testimonial/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Approved(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	testimonials, _ := models.Testimonial{}.GetApproved()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 3:04 PM"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.Bag["PageTitle"] = "All Approved Testimonials"
	tmpl.Bag["testimonials"] = testimonials
	tmpl.Bag["type"] = "Approved"

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/testimonial/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}
