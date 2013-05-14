package blog

import (
	"../../helpers/plate"
	"../../models"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)

	posts, _ := models.Post{}.GetAll()

	tmpl.FuncMap["formatDateForURL"] = func(dt time.Time) string {
		tlayout := "1-02-2006"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 3:04 PM"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.FuncMap["showCommentsLink"] = func(c models.Comments) bool {
		return len(c.Approved) > 0 && len(c.Unapproved) > 0
	}
	tmpl.Bag["PageTitle"] = "Blog Posts"
	tmpl.Bag["posts"] = posts

	tmpl.ParseFile("templates/blog/navigation.html", false)
	tmpl.ParseFile("templates/blog/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Categories(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)

	categories, _ := models.BlogCategory{}.GetAll()

	tmpl.Bag["PageTitle"] = "Blog Categories"
	tmpl.Bag["categories"] = categories

	tmpl.ParseFile("templates/blog/navigation.html", false)
	tmpl.ParseFile("templates/blog/categories.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func AddCategory(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	error, _ := url.QueryUnescape(r.URL.Query().Get("error"))

	if strings.TrimSpace(error) == "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Add Blog Category"
	tmpl.Bag["Type"] = "Add"
	tmpl.Bag["category"] = models.BlogCategory{}

	tmpl.ParseFile("templates/blog/navigation.html", false)
	tmpl.ParseFile("templates/blog/categoryform.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func EditCategory(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	error, _ := url.QueryUnescape(r.URL.Query().Get("error"))

	if strings.TrimSpace(error) == "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Edit Blog Category"
	tmpl.Bag["Type"] = "Edit"
	tmpl.Bag["category"], _ = models.BlogCategory{ID: id}.Get()

	tmpl.ParseFile("templates/blog/navigation.html", false)
	tmpl.ParseFile("templates/blog/categoryform.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func SaveCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	cat := models.BlogCategory{
		ID:   id,
		Name: r.FormValue("name"),
	}
	err := cat.Save()
	if err != nil {
		if cat.ID > 0 {
			http.Redirect(w, r, "/Blog/EditCategory/"+strconv.Itoa(cat.ID)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
		} else {
			http.Redirect(w, r, "/Blog/AddCategory?error="+url.QueryEscape(err.Error()), http.StatusFound)
		}
	}
	http.Redirect(w, r, "/Blog/Categories", http.StatusFound)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	cat := models.BlogCategory{ID: id}
	success := cat.Delete()
	plate.ServeFormatted(w, r, success)
}
