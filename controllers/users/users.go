package users

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

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	params := r.URL.Query()
	error := params.Get("error")

	u := models.User{ID: 0, Modules: make([]models.Module, 0)}
	modules, err := models.GetAllModules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.FuncMap["hasModule"] = func(modID int) bool {
		for _, m := range u.Modules {
			if modID == m.ID {
				return true
			}
		}
		return false
	}

	tmpl.Bag["u"] = u
	tmpl.Bag["error"] = error
	tmpl.Bag["modules"] = modules
	tmpl.ParseFile("templates/users/form.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	params := r.URL.Query()
	error := params.Get("error")
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := models.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	modules, err := models.GetAllModules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.FuncMap["isZero"] = func(num int) bool {
		return num == 0
	}
	tmpl.FuncMap["hasModule"] = func(modID int) bool {
		for _, m := range u.Modules {
			if modID == m.ID {
				return true
			}
		}
		return false
	}

	tmpl.Bag["error"] = error
	tmpl.Bag["modules"] = modules
	tmpl.Bag["u"] = u
	tmpl.ParseFile("templates/users/form.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}

}

func Save(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	active, _ := strconv.ParseBool(r.FormValue("isActive"))
	super, _ := strconv.ParseBool(r.FormValue("superUser"))

	u := models.User{
		ID:        id,
		Username:  r.FormValue("username"),
		Email:     r.FormValue("email"),
		Fname:     r.FormValue("fname"),
		Lname:     r.FormValue("lname"),
		Biography: r.FormValue("biography"),
		Photo:     r.FormValue("photo"),
		IsActive:  active,
		SuperUser: super,
	}

	pw1 := r.FormValue("password1")
	pw2 := r.FormValue("password2")
	if ValidatePassword(pw1, pw2) {
		epw, _ := models.Md5Encrypt(pw1)
		u.Password = epw
	}

	if id == 0 && u.Password == "" {
		http.Redirect(w, r, "/Users/Add?error=Password is Required.", http.StatusFound)
		return
	}

	err = u.Save()
	if err != nil {
		if u.ID == 0 {
			http.Redirect(w, r, "/Users/Add?error="+err.Error(), http.StatusFound)
		} else {
			http.Redirect(w, r, "/Users/Edit/"+strconv.Itoa(id)+"?error="+err.Error(), http.StatusFound)
		}
	}

	modules := r.Form["module"]
	u.SaveModules(modules)
	http.Redirect(w, r, "/Users/Edit/"+strconv.Itoa(u.ID), http.StatusFound)
}

func ValidatePassword(pw1 string, pw2 string) bool {
	valid := false
	if len(pw1) >= 8 && pw1 == pw2 {
		valid = true
	}

	return valid
}

func Delete(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	success := false
	id, err := strconv.Atoi(params.Get(":id"))
	if err == nil {
		u := models.User{}
		u.ID = id
		err = u.Delete()
		if err == nil {
			success = true
		}
	}
	successobj := struct {
		Success bool
	}{
		success,
	}
	plate.ServeFormatted(w, r, successobj)
}

func SetUserStatus(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	success := false
	id, err := strconv.Atoi(params.Get(":id"))
	if err == nil {
		u, err := models.GetUserByID(id)
		if err == nil {
			err = u.SetStatus(!u.IsActive)
			if err == nil {
				success = true
			}
		}
	}
	successobj := struct {
		Success bool
	}{
		success,
	}
	plate.ServeFormatted(w, r, successobj)
}

func MyAccount(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)

	params := r.URL.Query()
	error, _ := url.QueryUnescape(params.Get("error"))
	tmpl.Bag["error"] = error
	tmpl.ParseFile("templates/account/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	u, _ := tmpl.Bag["user"].(models.User)
	change := false

	email := r.FormValue("email")
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	bio := r.FormValue("biography")
	photo := r.FormValue("photo")
	password := r.FormValue("password1")
	confirm := r.FormValue("password2")

	if u.Email != strings.TrimSpace(email) {
		u.Email = strings.TrimSpace(email)
		change = true
	}

	if u.Fname != strings.TrimSpace(fname) {
		u.Fname = strings.TrimSpace(fname)
		change = true
	}

	if u.Lname != strings.TrimSpace(lname) {
		u.Lname = strings.TrimSpace(lname)
		change = true
	}

	if u.Biography != strings.TrimSpace(bio) {
		u.Biography = strings.TrimSpace(bio)
		change = true
	}

	if u.Photo != strings.TrimSpace(photo) {
		u.Photo = strings.TrimSpace(photo)
		change = true
	}

	if password != "" && (password != confirm || len(password) < 8) {
		http.Redirect(w, r, "/Account?error="+url.QueryEscape("Password does not match or does not meet requirements"), http.StatusFound)
		return
	} else {
		epw, _ := models.Md5Encrypt(password)
		u.Password = epw
		change = true
	}

	if change {
		err := u.Save()
		if err != nil {
			http.Redirect(w, r, "/Account?error="+url.QueryEscape(err.Error()), http.StatusFound)
			return
		}
	}

	http.Redirect(w, r, "/Account?error="+url.QueryEscape("Account Info Updated"), http.StatusFound)
	return
}
