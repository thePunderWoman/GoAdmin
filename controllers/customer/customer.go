package customer

import (
	"../../helpers/plate"
	"../../models"
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	error, _ := url.QueryUnescape(r.URL.Query().Get("error"))

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
	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
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

func Add(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)
	error, _ := url.QueryUnescape(r.URL.Query().Get("error"))

	customer := models.Customer{}
	session, _ := store.Get(r, "adminstuffs")
	if cjson := session.Flashes("customer"); len(cjson) > 0 {
		json.Unmarshal([]byte(cjson[0].(string)), &customer)
		session.Save(r, w)
	}
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
	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
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

func Save(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("cust_id"))
	stateID, _ := strconv.Atoi(r.FormValue("state"))
	dealertypeID, _ := strconv.Atoi(r.FormValue("dealer_type"))
	customerID, _ := strconv.Atoi(r.FormValue("customerID"))
	parentID, _ := strconv.Atoi(r.FormValue("parentID"))
	mapicscodeID, _ := strconv.Atoi(r.FormValue("mapixCodeID"))
	salesrepID, _ := strconv.Atoi(r.FormValue("salesRepID"))
	tierID, _ := strconv.Atoi(r.FormValue("tier"))
	dummy, _ := strconv.ParseBool(r.FormValue("isDummy"))
	show, _ := strconv.ParseBool(r.FormValue("showWebsite"))
	cust := models.Customer{
		ID:            id,
		Name:          strings.TrimSpace(r.FormValue("name")),
		Email:         strings.TrimSpace(r.FormValue("email")),
		Address:       strings.TrimSpace(r.FormValue("address")),
		Address2:      strings.TrimSpace(r.FormValue("address2")),
		City:          strings.TrimSpace(r.FormValue("city")),
		StateID:       stateID,
		PostalCode:    strings.TrimSpace(r.FormValue("postal_code")),
		Phone:         strings.TrimSpace(r.FormValue("phone")),
		Fax:           strings.TrimSpace(r.FormValue("fax")),
		ContactPerson: strings.TrimSpace(r.FormValue("contact_person")),
		DealerTypeID:  dealertypeID,
		Website:       strings.TrimSpace(r.FormValue("website")),
		CustomerID:    customerID,
		IsDummy:       dummy,
		ParentID:      parentID,
		SearchURL:     strings.TrimSpace(r.FormValue("searchURL")),
		ELocalURL:     strings.TrimSpace(r.FormValue("eLocalURL")),
		Logo:          strings.TrimSpace(r.FormValue("logo")),
		MapicsCodeID:  mapicscodeID,
		SalesRepID:    salesrepID,
		Tier:          tierID,
		ShowWebsite:   show,
	}
	if cust.Name == "" || cust.DealerTypeID == 0 {
		err := errors.New("Name and Dealer Type are required.")
		if cust.ID > 0 {
			http.Redirect(w, r, "/Customers/Edit/"+strconv.Itoa(cust.ID)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
			return
		} else {
			cjson, _ := json.Marshal(&cust)
			session, _ := store.Get(r, "adminstuffs")
			session.AddFlash(string(cjson), "customer")
			session.Save(r, w)
			http.Redirect(w, r, "/Customers/Add?error="+url.QueryEscape(err.Error()), http.StatusFound)
			return
		}
	}
	err := cust.Save()
	if err != nil {
		if cust.ID > 0 {
			http.Redirect(w, r, "/Customers/Edit/"+strconv.Itoa(cust.ID)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
			return
		} else {
			cjson, _ := json.Marshal(&cust)
			session, _ := store.Get(r, "adminstuffs")
			session.AddFlash(string(cjson), "customer")
			session.Save(r, w)
			http.Redirect(w, r, "/Customers/Add?error="+url.QueryEscape(err.Error()), http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, "/Customers/Edit/"+strconv.Itoa(cust.ID), http.StatusFound)
}

func Locations(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))

	customer, _ := models.Customer{ID: id}.Get()
	locations, _ := models.CustomerLocation{CustomerID: customer.ID}.GetAll()

	tmpl.Bag["PageTitle"] = "Customer Locations"
	tmpl.Bag["customer"] = customer
	tmpl.Bag["locations"] = locations

	tmpl.ParseFile("templates/customer/navigation.html", false)
	tmpl.ParseFile("templates/customer/locations.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func LocationsJSON(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	locations, _ := models.CustomerLocation{CustomerID: id}.GetAll()
	plate.ServeFormatted(w, r, locations)
}

func MapIconJSON(w http.ResponseWriter, r *http.Request) {
	icons, _ := models.MapIcon{}.GetAll()
	plate.ServeFormatted(w, r, icons)
}
