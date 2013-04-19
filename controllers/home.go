package controllers

import (
	//	"fmt"
	"../helpers/plate"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)

	tmpl.ParseFile("templates/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}
