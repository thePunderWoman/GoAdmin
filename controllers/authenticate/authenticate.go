package authenticate

import (
	"../../helpers/globals"
	"../../helpers/plate"
	"../../models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) bool {
	cook, err := r.Cookie("userID")
	userID := 0
	if err == nil && cook != nil {
		userID, err = strconv.Atoi(cook.Value)
		if err != nil {
			userID = 0
		}
	}
	if userID == 0 {
		http.Redirect(w, r, "/Authenticate", http.StatusFound)
	}

	return true
}

func Index(w http.ResponseWriter, r *http.Request) {
	var err error
	var tmpl plate.Template

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

	templates := append(globals.StandardLayout, "templates/auth/login.html")

	tmpl.DisplayMultiple(templates)
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := models.Authenticate(username, password)

	if err != nil {
		urlpath := "/authenticate/" + url.QueryEscape("Failed to log you into the system")
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
