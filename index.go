package main

import (
	"./controllers"
	"./controllers/authenticate"
	"./controllers/base"
	"./controllers/users"
	"./controllers/website"
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
	server.AddFilter(base.Base)

	// Authentication Routes
	server.Get("/Authenticate", authenticate.Index).NoFilter()
	server.Post("/Authenticate", authenticate.Login).NoFilter()
	server.Get("/Forgot", authenticate.Forgot).NoFilter()
	server.Post("/Forgot", authenticate.NewPassword).NoFilter()
	server.Get("/Signup", authenticate.SignUp).NoFilter()
	server.Post("/Signup", authenticate.Register).NoFilter()
	server.Get("/Logout", authenticate.Logout)
	server.Get("/Logout", authenticate.Logout)
	server.Get("/Account", users.MyAccount)
	server.Post("/Account", users.UpdateAccount)

	// User Routes
	server.Get("/Users", users.Index)
	server.Get("/Users/Add", users.Add)
	server.Get("/Users/Edit/:id", users.Edit)
	server.Post("/Users/Save/:id", users.Save)
	server.Get("/Users/SetUserStatus/:id", users.SetUserStatus)

	// Website Routes
	server.Get("/Website", website.Index)
	server.Get("/Website/Menus", website.Menus)
	server.Get("/Website/Menu/Add", website.Add)
	server.Get("/Website/Menu/SetPrimary/:id", website.SetPrimaryMenu)
	server.Post("/Website/Menu/Save/:id", website.Save)
	server.Get("/Website/Menu/Edit/:id", website.Edit)
	server.Get("/Website/Menu/:id", website.Menu)
	server.Post("/Website/Menu/Remove/:id", website.Remove)
	server.Post("/Website/Menu/Sort", website.MenuSort)
	server.Post("/Website/AddContentToMenu", website.AddContentToMenu)
	server.Post("/Website/RemoveContentAjax/:id", website.RemoveContentAjax)
	server.Get("/Website/SetPrimaryContent/:id/:menuid", website.SetPrimaryContent)
	server.Get("/Website/Link/Add/:id", website.AddLink)
	server.Post("/Website/Link/Add/:id", website.SaveLink)
	server.Get("/Website/checkContent/:id", website.CheckContent)
	server.Post("/Website/DeleteContent/:id", website.DeleteContent)
	server.Post("/Website/Content/Add", website.SaveContent)
	server.Get("/Website/Content/Add/:id", website.AddContent)
	server.Get("/Website/Content/Edit/:id", website.EditContent)
	server.Get("/Website/Content/Edit/:id/:revid", website.EditContent)

	// Home page route
	server.Get("/", controllers.Index)

	dir, _ := os.Getwd()

	server.Static("/", dir+"/"+"static")

	http.Handle("/", server)

	log.Println("Server running on port " + *globals.ListenAddr)

	log.Fatal(http.ListenAndServe(*globals.ListenAddr, nil))

}
