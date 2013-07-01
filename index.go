package main

import (
	"./controllers"
	"./controllers/authenticate"
	"./controllers/base"
	"./controllers/blog"
	"./controllers/contact"
	"./controllers/customer"
	"./controllers/faq"
	"./controllers/landingpage"
	"./controllers/news"
	"./controllers/salesrep"
	"./controllers/testimonial"
	"./controllers/users"
	"./controllers/video"
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
	server.Get("/Website/Link/Edit/:id", website.EditLink)
	server.Post("/Website/Link/Add/:id", website.SaveNewLink)
	server.Post("/Website/Link/Edit/:id", website.SaveLink)
	server.Get("/Website/checkContent/:id", website.CheckContent)
	server.Post("/Website/DeleteContent/:id", website.DeleteContent)
	server.Post("/Website/Content/Add", website.SaveNewContent)
	server.Get("/Website/Content/Add/:id", website.AddContent)
	server.Get("/Website/Content/Edit/:id", website.EditContent)
	server.Get("/Website/Content/Edit/:id/:revid", website.EditContent)
	server.Post("/Website/Content/Edit/:id/:revid", website.SaveContent)
	server.Get("/Website/CopyRevision/:id", website.CopyRevision)
	server.Get("/Website/ActivateRevision/:id", website.ActivateRevision)
	server.Get("/Website/DeleteRevision/:id", website.DeleteRevision)

	// Contact Manager
	server.Get("/Contact", contact.Index)
	server.Get("/Contact/ViewContact/:id", contact.View)
	server.Get("/Contact/Receivers", contact.Receivers)
	server.Get("/Contact/Types", contact.Types)
	server.Get("/Contact/AddReceiver", contact.AddReceiver)
	server.Get("/Contact/EditReceiver/:id", contact.EditReceiver)
	server.Post("/Contact/SaveReceiver", contact.SaveReceiver)
	server.Get("/Contact/DeleteReceiver/:id", contact.DeleteReceiver)
	server.Get("/Contact/AddType", contact.AddType)
	server.Post("/Contact/SaveType", contact.SaveType)
	server.Get("/Contact/DeleteType/:id", contact.DeleteType)

	// FAQ
	server.Get("/FAQ", faq.Index)
	server.Get("/FAQ/Add", faq.Add)
	server.Get("/FAQ/Edit/:id", faq.Edit)
	server.Post("/FAQ/Save", faq.Save)
	server.Get("/FAQ/Delete/:id", faq.Delete)

	// News
	server.Get("/News", news.Index)
	server.Get("/News/Add", news.Add)
	server.Get("/News/Edit/:id", news.Edit)
	server.Post("/News/Save", news.Save)
	server.Get("/News/Delete/:id", news.Delete)

	// Video
	server.Get("/Video", video.Index)
	server.Post("/Video/UpdateSort", video.Sort)
	server.Post("/Video/Delete", video.Delete)
	server.Get("/Video/AddVideo", video.Add)

	// Testimonials
	server.Get("/Testimonial", testimonial.Index)
	server.Get("/Testimonial/Approved", testimonial.Approved)
	server.Get("/Testimonial/Remove", testimonial.Remove)
	server.Get("/Testimonial/SetApproval", testimonial.SetApproval)

	// Landing Pages
	server.Get("/LandingPages", landingpage.Index)
	server.Get("/LandingPages/Past", landingpage.Past)
	server.Get("/LandingPages/Add", landingpage.Add)
	server.Get("/LandingPages/Edit/:id", landingpage.Edit)
	server.Get("/LandingPages/AddData", landingpage.AddData)
	server.Get("/LandingPages/RemoveData/:id", landingpage.RemoveData)
	server.Post("/LandingPages/Add", landingpage.New)
	server.Post("/LandingPages/Save", landingpage.Save)
	server.Get("/LandingPages/Remove/:id", landingpage.Remove)

	// Sales Reps
	server.Get("/SalesRep", salesrep.Index)
	server.Get("/SalesRep/Add", salesrep.Add)
	server.Get("/SalesRep/Edit/:id", salesrep.Edit)
	server.Post("/SalesRep/Save", salesrep.Save)
	server.Get("/SalesRep/Delete/:id", salesrep.Delete)

	// Blog
	server.Get("/Blog", blog.Index)
	server.Get("/Blog/Categories", blog.Categories)
	server.Get("/Blog/AddCategory", blog.AddCategory)
	server.Get("/Blog/EditCategory/:id", blog.EditCategory)
	server.Post("/Blog/SaveCategory", blog.SaveCategory)
	server.Get("/Blog/DeleteCategory/:id", blog.DeleteCategory)
	server.Get("/Blog/Add", blog.AddPost)
	server.Get("/Blog/Edit/:id", blog.EditPost)
	server.Post("/Blog/Save", blog.SavePost)
	server.Get("/Blog/Delete/:id", blog.DeletePost)
	server.Get("/Blog/Comments", blog.Comments)
	server.Get("/Blog/Comment/:id", blog.Comment)
	server.Get("/Blog/ApproveComment/:id", blog.ApproveComment)
	server.Get("/Blog/DeleteComment/:id", blog.DeleteComment)
	server.Get("/Blog/PostComments/:id", blog.PostComments)

	// Customers
	server.Get("/Customers", customer.Index)
	server.Get("/Customers/Edit/:id", customer.Edit)
	server.Get("/Customers/Add", customer.Add)
	server.Post("/Customers/Save", customer.Save)

	// Customer Locations
	server.Get("/Customers/Locations/:id", customer.Locations)
	server.Get("/Customers/LocationJSON/:id", customer.LocationsJSON)
	server.Get("/Customers/MapIcons", customer.MapIconJSON)
	server.Get("/Customers/AddLocation/:id", customer.AddLocation)
	server.Get("/Customers/EditLocation/:id", customer.EditLocation)
	server.Post("/Customers/Location/Save", customer.SaveLocation)
	server.Get("/Customers/Location/Delete/:id", customer.DeleteLocation)
	server.Get("/Customers/Location/Delete/:id", customer.DeleteLocation)
	server.Get("/Customers/PopulateLocations/:id", customer.PopulateCustomerLocations)

	// Customer Users
	server.Get("/Customers/CustomerUsers/:id", customer.CustomerUsers)
	server.Get("/Customers/Users", customer.AllCustomerUsers)

	// Home page route
	server.Get("/", controllers.Index)

	dir, _ := os.Getwd()

	server.Static("/", dir+"/"+"static")

	http.Handle("/", server)

	log.Println("Server running on port " + *globals.ListenAddr)

	log.Fatal(http.ListenAndServe(*globals.ListenAddr, nil))

}
