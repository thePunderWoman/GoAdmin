package authenticate

import (
	"../../helpers/plate"
	"../../models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	TemplateFiles = []string{
		"templates/shared/head.html",
		"templates/shared/header.html",
		"templates/shared/navigation.html",
		"templates/shared/footer.html",
	}
)

func Index(w http.ResponseWriter, r *http.Request) {
	var err error
	var tmpl *plate.Template

	params := r.URL.Query()
	error := params.Get(":error")
	error, _ = url.QueryUnescape(error)
	server := plate.NewServer()

	tmpl, err = server.Template(w)

	if err != nil {
		plate.Serve404(w, err.Error())
		return
	}

	tmpl.Bag["Error"] = strings.ToTitle(error)
	tmpl.Bag["CurrentYear"] = time.Now().Year()
	tmpl.Bag["userID"] = 0

	tmpl.FuncMap["isNotNull"] = func(str string) bool {
		if strings.TrimSpace(str) != "" && len(strings.TrimSpace(str)) > 0 {
			return true
		}
		return false
	}
	tmpl.FuncMap["isLoggedIn"] = func() bool {
		return false
	}

	templates := append(TemplateFiles, "templates/auth/login.html")

	tmpl.DisplayMultiple(templates)
}

func Forgot(w http.ResponseWriter, r *http.Request) {
	var err error
	var tmpl *plate.Template

	params := r.URL.Query()
	error := params.Get("error")
	error, _ = url.QueryUnescape(error)
	server := plate.NewServer()

	tmpl, err = server.Template(w)

	if err != nil {
		plate.Serve404(w, err.Error())
		return
	}

	tmpl.Bag["Error"] = strings.ToTitle(error)
	tmpl.Bag["CurrentYear"] = time.Now().Year()
	tmpl.Bag["userID"] = 0

	tmpl.FuncMap["isNotNull"] = func(str string) bool {
		if strings.TrimSpace(str) != "" && len(strings.TrimSpace(str)) > 0 {
			return true
		}
		return false
	}
	tmpl.FuncMap["isLoggedIn"] = func() bool {
		return false
	}

	templates := append(TemplateFiles, "templates/auth/forgot.html")

	tmpl.DisplayMultiple(templates)
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := models.Authenticate(username, password)

	if err != nil {
		urlpath := "/authenticate?error=" + url.QueryEscape("Failed to log you into the system")
		http.Redirect(w, r, urlpath, http.StatusFound)
	} else {
		cook := http.Cookie{
			Name:    "userID",
			Value:   strconv.Itoa(user.ID),
			Expires: time.Now().AddDate(2, 0, 0),
		}

		cook2 := http.Cookie{
			Name:    "username",
			Value:   user.Username,
			Expires: time.Now().AddDate(2, 0, 0),
		}

		cook3 := http.Cookie{
			Name:    "superUser",
			Value:   strconv.FormatBool(user.SuperUser),
			Expires: time.Now().AddDate(2, 0, 0),
		}

		cook4 := http.Cookie{
			Name:    "name",
			Value:   user.Fname + " " + user.Lname,
			Expires: time.Now().AddDate(2, 0, 0),
		}

		http.SetCookie(w, &cook)
		http.SetCookie(w, &cook2)
		http.SetCookie(w, &cook3)
		http.SetCookie(w, &cook4)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// expire cookie
	cook, err := r.Cookie("userID")

	if err == nil {
		cook.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, cook)
	}

	cook2, err := r.Cookie("username")
	if err == nil {
		cook2.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, cook2)
	}

	cook3, err := r.Cookie("superUser")
	if err == nil {
		cook3.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, cook3)
	}

	cook4, err := r.Cookie("name")
	if err == nil {
		cook4.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, cook4)
	}

	http.Redirect(w, r, "/authenticate", http.StatusFound)
}

func NewPassword(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	var user models.User
	var err error
	if strings.TrimSpace(username) != "" {
		user, err = models.GetUserByUsername(username)
		if err != nil {
			http.Redirect(w, r, "/Forgot?error="+url.QueryEscape("The username you specified was not found"), http.StatusFound)
			return
		}
	} else if strings.TrimSpace(email) != "" {
		user, err = models.GetUserByEmail(email)
		if err != nil {
			http.Redirect(w, r, "/Forgot?error="+url.QueryEscape("The email address you specified was not found"), http.StatusFound)
			return
		}
	} else {
		http.Redirect(w, r, "/Forgot?error="+url.QueryEscape("No Email or Username was entered."), http.StatusFound)
		return
	}
	user.ResetPassword()
	http.Redirect(w, r, "/Forgot?error="+url.QueryEscape("Reset successful. Check your email for your new password."), http.StatusFound)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var err error
	var tmpl *plate.Template

	params := r.URL.Query()
	error := params.Get("error")
	error, _ = url.QueryUnescape(error)

	fname := params.Get("fname")
	fname, _ = url.QueryUnescape(fname)

	lname := params.Get("lname")
	lname, _ = url.QueryUnescape(lname)

	email := params.Get("email")
	email, _ = url.QueryUnescape(email)

	username := params.Get("username")
	username, _ = url.QueryUnescape(username)

	var submitted bool
	submitted, err = strconv.ParseBool(params.Get("submitted"))
	if err != nil {
		submitted = false
	}

	server := plate.NewServer()

	tmpl, err = server.Template(w)

	if err != nil {
		plate.Serve404(w, err.Error())
		return
	}

	tmpl.Bag["Error"] = strings.ToTitle(error)
	tmpl.Bag["Fname"] = strings.TrimSpace(fname)
	tmpl.Bag["Lname"] = strings.TrimSpace(lname)
	tmpl.Bag["Email"] = strings.TrimSpace(email)
	tmpl.Bag["Username"] = strings.TrimSpace(username)
	tmpl.Bag["CurrentYear"] = time.Now().Year()
	tmpl.Bag["Submitted"] = submitted
	tmpl.Bag["userID"] = 0

	tmpl.FuncMap["isNotNull"] = func(str string) bool {
		if strings.TrimSpace(str) != "" && len(strings.TrimSpace(str)) > 0 {
			return true
		}
		return false
	}
	tmpl.FuncMap["isLoggedIn"] = func() bool {
		return false
	}

	templates := append(TemplateFiles, "templates/auth/signup.html")

	tmpl.DisplayMultiple(templates)

}

func Register(w http.ResponseWriter, r *http.Request) {
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	email := r.FormValue("email")
	username := r.FormValue("username")
	if strings.TrimSpace(fname) == "" || strings.TrimSpace(lname) == "" || strings.TrimSpace(email) == "" || strings.TrimSpace(username) == "" {
		qvals := url.Values{}
		qvals.Add("error", "Please fill out required fields")
		qvals.Add("fname", fname)
		qvals.Add("lname", lname)
		qvals.Add("email", email)
		qvals.Add("username", username)
		http.Redirect(w, r, "/Signup?"+qvals.Encode(), http.StatusFound)
		return
	}

	user := models.User{
		Username: username,
		Email:    email,
		Fname:    fname,
		Lname:    lname,
	}

	err := user.Save()
	if err != nil {
		qvals := url.Values{}
		qvals.Add("error", err.Error())
		qvals.Add("fname", fname)
		qvals.Add("lname", lname)
		qvals.Add("email", email)
		qvals.Add("username", username)
		http.Redirect(w, r, "/Signup?"+qvals.Encode(), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/Signup?submitted=true", http.StatusFound)
}
