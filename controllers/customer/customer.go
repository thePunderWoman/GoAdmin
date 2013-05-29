package customer

import (
	"../../helpers/plate"
	"../../models"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	//"net/url"
	"strconv"
	//"strings"
)

var store = sessions.NewCookieStore([]byte("adminstuffs"))

func Index(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)

	customers, _ := models.Customer{}.GetAll()

	tmpl.Bag["PageTitle"] = "View Customers"
	tmpl.Bag["customers"] = customers

	tmpl.ParseFile("templates/customer/navigation.html", false)
	tmpl.ParseFile("templates/customer/index.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))

	customer, _ := models.Customer{ID: id}.Get()
	var customers models.Customers
	var dealertypes models.DealerTypes
	var dealertiers models.DealerTiers
	var salesreps models.SalesReps
	var mapicscodes models.MapicsCodes
	var countries models.Countries

	custchan := make(chan int)
	typechan := make(chan int)
	tierchan := make(chan int)
	repchan := make(chan int)
	codechan := make(chan int)
	statechan := make(chan int)

	go func(ch chan int) {
		dealertypes, _ = models.DealerType{}.GetAll()
		ch <- 1
	}(typechan)

	go func(ch chan int) {
		dealertiers, _ = models.DealerTier{}.GetAll()
		ch <- 1
	}(tierchan)

	go func(ch chan int) {
		salesreps, _ = models.SalesRep{}.GetAll()
		ch <- 1
	}(repchan)

	go func(ch chan int) {
		mapicscodes, _ = models.MapicsCode{}.GetAll()
		ch <- 1
	}(codechan)

	go func(ch chan int) {
		countries, _ = models.Country{}.GetAll()
		ch <- 1
	}(statechan)

	go func(ch chan int) {
		customers, _ = models.Customer{}.GetAllSimple()
		ch <- 1
	}(custchan)

	<-typechan
	<-tierchan
	<-repchan
	<-codechan
	<-statechan
	<-custchan

	tmpl.FuncMap["equals"] = func(id int, cid int) bool {
		return id == cid
	}
	tmpl.Bag["PageTitle"] = "Edit Customer"
	tmpl.Bag["customer"] = customer
	tmpl.Bag["customers"] = customers
	tmpl.Bag["dealertypes"] = dealertypes
	tmpl.Bag["tiers"] = dealertiers
	tmpl.Bag["countries"] = countries
	tmpl.Bag["salesreps"] = salesreps
	tmpl.Bag["mapicscodes"] = mapicscodes

	tmpl.ParseFile("templates/customer/navigation.html", false)
	tmpl.ParseFile("templates/customer/customerform.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}
