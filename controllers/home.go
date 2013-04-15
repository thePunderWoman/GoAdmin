package controllers

import (
	//	"fmt"
	"../helpers/globals"
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

	tmpl.Layout = "layout.html"
	templates := append(globals.StandardLayout, "templates/index.html")

	tmpl.DisplayMultiple(templates)
}
