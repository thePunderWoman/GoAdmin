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
	"time"
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

func MassUpload(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	tmpl.Bag["PageTitle"] = "Mass Upload Customers"
	tmpl.ParseFile("templates/customer/navigation.html", false)
	tmpl.ParseFile("templates/customer/massupload.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
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

func AddLocation(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	error, _ := url.QueryUnescape(r.URL.Query().Get("error"))

	custchan := make(chan int)
	countrychan := make(chan int)
	customer := models.Customer{}
	countries := models.Countries{}
	location := models.CustomerLocation{CustomerID: customer.ID}

	session, _ := store.Get(r, "adminstuffs")
	if cjson := session.Flashes("location"); len(cjson) > 0 {
		json.Unmarshal([]byte(cjson[0].(string)), &location)
		session.Save(r, w)
	}

	go func(ch chan int) {
		customer, _ = models.Customer{ID: id}.Get()
		ch <- 1
	}(custchan)
	go func(ch chan int) {
		countries, _ = models.Country{}.GetAll()
		ch <- 1
	}(countrychan)
	<-custchan
	<-countrychan

	tmpl.FuncMap["equals"] = func(locstateID int, stateID int) bool {
		return locstateID == stateID
	}
	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Customer Locations"
	tmpl.Bag["customer"] = customer
	tmpl.Bag["countries"] = countries
	tmpl.Bag["location"] = location

	tmpl.ParseFile("templates/customer/navigation.html", false)
	tmpl.ParseFile("templates/customer/locationform.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func EditLocation(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	error, _ := url.QueryUnescape(r.URL.Query().Get("error"))

	custchan := make(chan int)
	countrychan := make(chan int)
	customer := models.Customer{}
	countries := models.Countries{}
	location := models.CustomerLocation{}

	location, _ = models.CustomerLocation{ID: id}.Get()
	go func(ch chan int) {
		customer, _ = models.Customer{ID: location.CustomerID}.Get()
		ch <- 1
	}(custchan)
	go func(ch chan int) {
		countries, _ = models.Country{}.GetAll()
		ch <- 1
	}(countrychan)

	<-custchan
	<-countrychan

	tmpl.FuncMap["equals"] = func(locstateID int, stateID int) bool {
		return locstateID == stateID
	}
	if strings.TrimSpace(error) != "" {
		tmpl.Bag["error"] = error
	}
	tmpl.Bag["PageTitle"] = "Customer Locations"
	tmpl.Bag["customer"] = customer
	tmpl.Bag["countries"] = countries
	tmpl.Bag["location"] = location

	tmpl.ParseFile("templates/customer/navigation.html", false)
	tmpl.ParseFile("templates/customer/locationform.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func SaveLocation(w http.ResponseWriter, r *http.Request) {
	customerID, _ := strconv.Atoi(r.FormValue("customerID"))
	locationID, _ := strconv.Atoi(r.FormValue("locationID"))
	stateID, _ := strconv.Atoi(r.FormValue("state"))
	var latitude, longitude float64
	latitude, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		latitude = 0
	}
	longitude, err = strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		longitude = 0
	}
	location := models.CustomerLocation{
		ID:         locationID,
		CustomerID: customerID,
		Name:       strings.TrimSpace(r.FormValue("name")),
		Address:    strings.TrimSpace(r.FormValue("address")),
		City:       strings.TrimSpace(r.FormValue("city")),
		StateID:    stateID,
		PostalCode: strings.TrimSpace(r.FormValue("postalCode")),
		Email:      strings.TrimSpace(r.FormValue("email")),
		Phone:      strings.TrimSpace(r.FormValue("phone")),
		Fax:        strings.TrimSpace(r.FormValue("fax")),
		Latitude:   latitude,
		Longitude:  longitude,
	}
	err = location.Save()
	if err != nil {
		if location.ID > 0 {
			http.Redirect(w, r, "/Customers/EditLocation/"+strconv.Itoa(location.ID)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
			return
		} else {
			cjson, _ := json.Marshal(&location)
			session, _ := store.Get(r, "adminstuffs")
			session.AddFlash(string(cjson), "location")
			session.Save(r, w)
			http.Redirect(w, r, "/Customers/AddLocation/"+strconv.Itoa(customerID)+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, "/Customers/Locations/"+strconv.Itoa(customerID), http.StatusFound)
}

func DeleteLocation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	success := models.CustomerLocation{ID: id}.Delete()
	plate.ServeFormatted(w, r, success)
}

func PopulateCustomerLocations(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if id > 0 {
		locations, err := models.CustomerLocation{CustomerID: id}.GetAllNoGeo()
		if err == nil && len(locations) > 0 {
			for _, loc := range locations {
				go func(location models.CustomerLocation) {
					location.Save()
				}(loc)
			}
		}
	}
	http.Redirect(w, r, "/Customers/Locations/"+strconv.Itoa(id), http.StatusFound)
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

func CustomerUsers(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))

	customer, _ := models.Customer{ID: id}.Get()
	users, _ := models.CustomerUser{CustID: customer.ID}.GetAllByCustomer()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/06 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}

	tmpl.Bag["PageTitle"] = "Customer Users"
	tmpl.Bag["customer"] = customer
	tmpl.Bag["users"] = users

	tmpl.ParseFile("templates/customer/navigation.html", false)
	tmpl.ParseFile("templates/customer/users.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func AllCustomerUsers(w http.ResponseWriter, r *http.Request) {

	tmpl := plate.NewTemplate(w)

	users, _ := models.CustomerUser{}.GetAll()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/06 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}

	tmpl.Bag["PageTitle"] = "All Customer Users"
	tmpl.Bag["users"] = users

	tmpl.ParseFile("templates/customer/navigation.html", false)
	tmpl.ParseFile("templates/customer/allusers.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func EditCustomerUser(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	id := r.URL.Query().Get(":id")

	custuser, _ := models.CustomerUser{ID: id}.Get()
	customer, _ := models.Customer{ID: custuser.CustID}.Get()

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/06 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}

	tmpl.Bag["PageTitle"] = "Edit Customer User"
	tmpl.Bag["customer"] = customer
	tmpl.Bag["custuser"] = custuser

	tmpl.ParseFile("templates/customer/navigation.html", false)
	tmpl.ParseFile("templates/customer/userform.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func AddCustomerUser(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))

	customer, _ := models.Customer{ID: id}.Get()
	custuser := models.CustomerUser{
		CustID:     customer.ID,
		CustomerID: customer.CustomerID,
	}

	tmpl.FuncMap["formatDate"] = func(dt time.Time) string {
		tlayout := "01/02/06 3:04PM MST"
		Local, _ := time.LoadLocation("US/Central")
		return dt.In(Local).Format(tlayout)
	}

	tmpl.Bag["PageTitle"] = "Add Customer User"
	tmpl.Bag["customer"] = customer
	tmpl.Bag["custuser"] = custuser

	tmpl.ParseFile("templates/customer/navigation.html", false)
	tmpl.ParseFile("templates/customer/userform.html", false)

	err := tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}
