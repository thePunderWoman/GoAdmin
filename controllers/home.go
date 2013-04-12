package controllers

import (
	//	"fmt"
	"../helpers/globals"
	"../helpers/plate"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	server := plate.NewServer()

	tmpl, _ := server.Template(w, r)

	templates := append(globals.StandardLayout, "templates/index.html")

	tmpl.DisplayMultiple(templates)
}
