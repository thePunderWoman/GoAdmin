package controllers

import (
	//	"fmt"
	"../helpers/plate"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {

	server := plate.NewServer()

	tmpl, err := plate.GetTemplate()
	if err != nil {
		tmpl, err = server.Template(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tmpl.ParseFile("templates/index.html", false)

	tmpl.Display(w)
}
