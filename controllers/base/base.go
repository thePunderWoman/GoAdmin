package base

import (
	"../../helpers/plate"
	"net/http"
	"os"
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
	}

	server := plate.NewServer()
	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := plate.GetTemplate()
	if err != nil {
		tmpl, err = server.Template(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tmpl.Bag["CurrentYear"] = time.Now().Year()
	tmpl.Bag["userID"] = userID

	tmpl.FuncMap = template.FuncMap{
		"isNotNull": func(str string) bool {
			if strings.TrimSpace(str) != "" && len(strings.TrimSpace(str)) > 0 {
				return true
			}
			return false
		},
	}

	for _, file := range TemplateFiles {
		err = tmpl.ParseFile(file)
	}

	plate.SetTemplate(tmpl)
}
