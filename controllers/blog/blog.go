package blog

import (
	"../../helpers/plate"
	"../../models"
	"encoding/json"
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

func AddPost(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	error, _ := url.QueryUnescape(r.URL.Query().Get("error"))
	message, _ := url.QueryUnescape(r.URL.Query().Get("message"))
	post := models.Post{}
	session, _ := store.Get(r, "adminstuffs")
	if pjson := session.Flashes("post"); len(pjson) > 0 {
		json.Unmarshal([]byte(pjson[0].(string)), &post)
		session.Save(r, w)
	}

	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}
	if strings.TrimSpace(message) != "" {
		tmpl.Bag["message"] = message
	}
	u := models.User{}

	tmpl.FuncMap["isUser"] = func(uid int) bool {
		return uid == post.UserID
	}
	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 3:04 PM"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.FuncMap["hasCategory"] = func(cid int) bool {
		for _, cat := range post.Categories {
			if cid == cat.ID {
				return true
			}
		}
		return false
	}
	tmpl.Bag["PageTitle"] = "Add Blog Post"
	tmpl.Bag["type"] = "Add"
	tmpl.Bag["categories"], _ = models.BlogCategory{}.GetAll()
	tmpl.Bag["users"], _ = u.GetAll()
	tmpl.Bag["post"] = post

	tmpl.ParseFile("templates/blog/navigation.html", false)
	tmpl.ParseFile("templates/blog/postform.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	error, _ := url.QueryUnescape(r.URL.Query().Get("error"))
	message, _ := url.QueryUnescape(r.URL.Query().Get("message"))
	log.Println(message)
	post, _ := models.Post{ID: id}.Get()
	session, _ := store.Get(r, "adminstuffs")
	if pjson := session.Flashes("post"); len(pjson) > 0 {
		json.Unmarshal([]byte(pjson[0].(string)), &post)
		session.Save(r, w)
	}

	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}
	if strings.TrimSpace(message) != "" {
		tmpl.Bag["message"] = message
	}
	u := models.User{}
	users, _ := u.GetAll()

	tmpl.FuncMap["isUser"] = func(uid int) bool {
		return uid == post.UserID
	}
	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 3:04 PM"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.FuncMap["hasCategory"] = func(cid int) bool {
		for _, cat := range post.Categories {
			if cid == cat.ID {
				return true
			}
		}
		return false
	}
	tmpl.Bag["PageTitle"] = "Edit Blog Post"
	tmpl.Bag["type"] = "Edit"
	tmpl.Bag["categories"], _ = models.BlogCategory{}.GetAll()
	tmpl.Bag["users"] = users
	tmpl.Bag["post"] = post

	tmpl.ParseFile("templates/blog/navigation.html", false)
	tmpl.ParseFile("templates/blog/postform.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func SavePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	userid, _ := strconv.Atoi(r.FormValue("userid"))
	publish, _ := strconv.ParseBool(r.FormValue("publish"))
	published := strings.TrimSpace(r.FormValue("published"))
	var categories models.BlogCategories
	categorylist := r.Form["categories"]
	for _, catid := range categorylist {
		cid, _ := strconv.Atoi(catid)
		cat := models.BlogCategory{
			ID:     cid,
			Active: true,
		}
		categories = append(categories, cat)
	}
	post := models.Post{
		ID:         id,
		Title:      r.FormValue("title"),
		Content:    r.FormValue("content"),
		UserID:     userid,
		MetaTitle:  r.FormValue("meta_title"),
		MetaDesc:   r.FormValue("meta_description"),
		Keywords:   r.FormValue("keywords"),
		Categories: categories,
	}
	err := post.Save()
	if err != nil {
		pjson, _ := json.Marshal(&post)
		session, _ := store.Get(r, "adminstuffs")
		session.AddFlash(string(pjson), "post")
		session.Save(r, w)
		if post.ID > 0 {
			http.Redirect(w, r, "/Blog/Edit/"+strconv.Itoa(post.ID)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
		} else {
			http.Redirect(w, r, "/Blog/Add?error="+url.QueryEscape(err.Error()), http.StatusFound)
			//http.Redirect(w, r, "/Blog/Add?message="+url.QueryEscape("Blog Post Saved successfully"), http.StatusFound)
		}
	}
	if publish && published == "" {
		post.Publish()
	} else if !publish && published != "" {
		post.UnPublish()
	}
	http.Redirect(w, r, "/Blog/Edit/"+strconv.Itoa(post.ID)+"?message="+url.QueryEscape("Blog Post Saved successfully"), http.StatusFound)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	post := models.Post{ID: id}
	success := post.Delete()
	plate.ServeFormatted(w, r, success)
}

func Comments(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	comments, _ := models.Comment{}.GetAll()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/2006 3:04 PM"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}
	tmpl.Bag["PageTitle"] = "Blog Comments"
	tmpl.Bag["comments"] = comments

	tmpl.ParseFile("templates/blog/navigation.html", false)
	tmpl.ParseFile("templates/blog/comments.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Comment(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	comment, _ := models.Comment{ID: id}.Get()
	post, _ := models.Post{ID: comment.PostID}.Get()

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
	tmpl.Bag["PageTitle"] = "View Blog Comment"
	tmpl.Bag["comment"] = comment
	tmpl.Bag["post"] = post

	tmpl.ParseFile("templates/blog/navigation.html", false)
	tmpl.ParseFile("templates/blog/comment.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func ApproveComment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	comment := models.Comment{ID: id}
	success := comment.Approve()
	plate.ServeFormatted(w, r, success)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	comment := models.Comment{ID: id}
	success := comment.Delete()
	plate.ServeFormatted(w, r, success)
}

func PostComments(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))

	post, _ := models.Post{ID: id}.Get()

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
	tmpl.Bag["PageTitle"] = "Post Comments"
	tmpl.Bag["post"] = post

	tmpl.ParseFile("templates/blog/navigation.html", false)
	tmpl.ParseFile("templates/blog/postcomments.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}
