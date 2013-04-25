package base

import (
	"../../helpers/plate"
	"../../models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	TemplateFiles = []string{
		"layout.html",
		"templates/shared/head.html",
		"templates/shared/header.html",
		"templates/shared/navigation.html",
		"templates/shared/footer.html",
	}
)

func Base(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	cook, err := r.Cookie("userID")
	userID := 0
	if err == nil && cook != nil {
		userID, err = strconv.Atoi(cook.Value)
		if err != nil {
			userID = 0
		}
	}
	if userID == 0 {
		// user is not logged in
		http.Redirect(w, r, "/Authenticate", http.StatusFound)
		return
	}
	user, err := models.GetUserByID(userID)
	if err != nil {
		// user is not logged in
		http.Redirect(w, r, "/Authenticate", http.StatusFound)
		return
	}
	if !user.HasModuleAccess(r.URL) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	server := plate.NewServer()
	tmpl, err := server.Template(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Bag["user"] = user
	tmpl.Bag["CurrentYear"] = time.Now().Year()
	tmpl.Bag["userID"] = userID

	tmpl.FuncMap["isNotNull"] = func(str string) bool {
		if strings.TrimSpace(str) != "" && len(strings.TrimSpace(str)) > 0 {
			return true
		}
		return false
	}
	tmpl.FuncMap["isZero"] = func(num int) bool {
		return num == 0
	}

	tmpl.FuncMap["isLoggedIn"] = func() bool {
		return userID > 0
	}

	for _, file := range TemplateFiles {
		err = tmpl.ParseFile(file, false)
		if err != nil {
			log.Println(err)
		}
	}

	plate.SetTemplate(tmpl)
}
