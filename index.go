package main

import (
	"./controllers/authenticate"
	"./controllers/base"
	"./helpers/database"
	"./helpers/globals"
	_ "./helpers/mimetypes"
	"./helpers/plate"
	"log"
	"net/http"
	"os"
)

var (
	CorsHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		return
	}
	AuthHandler = func(w http.ResponseWriter, r *http.Request) {
		authenticate.AuthHandler(w, r)
		return
	}
)

const (
	port = "80"
)

func main() {
	err := database.PrepareAll()
	if err != nil {
		log.Fatal(err)
	}

	globals.SetGlobals()
	server := plate.NewServer("doughboy")

	server.AddFilter(CorsHandler)

	server.Get("/Authenticate", authenticate.Index).NoFilter()
	server.Post("/Authenticate", authenticate.Login).NoFilter()
	server.Get("/Logout", authenticate.Logout)

	server.Get("/", controllers.Index).AddFilter(AuthHandler)

	dir, _ := os.Getwd()

	server.Static("/", dir+"/"+"static")

	http.Handle("/", server)

	log.Println("Server running on port " + *globals.ListenAddr)

	log.Fatal(http.ListenAndServe(*globals.ListenAddr, nil))

}
