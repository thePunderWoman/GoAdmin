package models

import (
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	"log"
	"time"
)

type Testimonial struct {
	ID        int
	Rating    float64
	Title     string
	Content   string
	DateAdded time.Time
	Approved  bool
	Active    bool
	FirstName string
	LastName  string
	Location  string
}

func (t Testimonial) GetUnapproved() ([]Testimonial, error) {
	testimonials := make([]Testimonial, 0)
	sel, err := database.GetStatement("GetAllTestimonialsStmt")
	if err != nil {
		return testimonials, err
	}
	sel.Bind(0)
	rows, res, err := sel.Exec()
	if err != nil {
		return testimonials, err
	}
	ch := make(chan Testimonial)
	for _, row := range rows {
		go t.PopulateTestimonial(row, res, ch)
	}
	for _, _ = range rows {
		testimonials = append(testimonials, <-ch)
	}
	return testimonials, nil
}

func (t Testimonial) GetApproved() ([]Testimonial, error) {
	testimonials := make([]Testimonial, 0)
	sel, err := database.GetStatement("GetAllTestimonialsStmt")
	if err != nil {
		return testimonials, err
	}
	sel.Bind(1)
	rows, res, err := sel.Exec()
	if err != nil {
		return testimonials, err
	}
	ch := make(chan Testimonial)
	for _, row := range rows {
		go t.PopulateTestimonial(row, res, ch)
	}
	for _, _ = range rows {
		testimonials = append(testimonials, <-ch)
	}
	return testimonials, nil
}

func (t *Testimonial) Percent() int {
	return int((t.Rating / 5) * 100)
}

func (t Testimonial) Remove() {
	upd, err := database.GetStatement("DeleteTestimonialStmt")
	if err != nil {
		log.Println(err)
		return
	}
	upd.Bind(t.ID)
	upd.Exec()
}

func (t Testimonial) SetApproval() bool {
	upd, err := database.GetStatement("SetTestimonialApprovalStmt")
	if err != nil {
		log.Println(err)
		return !t.Approved
	}
	upd.Bind(t.Approved, t.ID)
	upd.Exec()
	return t.Approved
}

func (t Testimonial) PopulateTestimonial(row mysql.Row, res mysql.Result, ch chan Testimonial) {
	testimonial := Testimonial{
		ID:        row.Int(res.Map("testimonialID")),
		Rating:    row.Float(res.Map("rating")),
		Title:     row.Str(res.Map("title")),
		Content:   row.Str(res.Map("testimonial")),
		DateAdded: row.Time(res.Map("dateAdded"), UTC),
		Approved:  row.Bool(res.Map("approved")),
		Active:    row.Bool(res.Map("active")),
		FirstName: row.Str(res.Map("first_name")),
		LastName:  row.Str(res.Map("last_name")),
		Location:  row.Str(res.Map("location")),
	}
	ch <- testimonial
}
