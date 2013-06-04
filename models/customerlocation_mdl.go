package models

import (
	//"../helpers/UDF"
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	//"sort"
)

type CustomerLocation struct {
	ID              int
	CustomerID      int
	Name            string
	Address         string
	City            string
	StateID         int
	State           State
	PostalCode      string
	Email           string
	Phone           string
	Fax             string
	Latitude        float64
	Longitude       float64
	ContactPerson   string
	IsPrimary       bool
	ShippingDefault bool
}

type CustomerLocations []CustomerLocation

func (c CustomerLocation) GetAll() (locations CustomerLocations, err error) {
	stateMap := make(map[int]State)
	statechan := make(chan int)

	go func(ch chan int) {
		states, _ := State{}.GetAll()
		stateMap = states.ToMap()
		ch <- 1
	}(statechan)

	sel, err := database.GetStatement("GetCustomerLocationsStmt")
	if err != nil {
		return locations, err
	}
	sel.Bind(c.CustomerID)
	rows, res, err := sel.Exec()
	if err != nil {
		return locations, err
	}
	<-statechan

	ch := make(chan CustomerLocation)
	for _, row := range rows {
		go c.PopulateLocation(row, res, ch)
	}
	for _, _ = range rows {
		loc := <-ch
		loc.State = stateMap[loc.StateID]
		locations = append(locations, loc)
	}
	return
}

func (c CustomerLocation) PopulateLocation(row mysql.Row, res mysql.Result, ch chan CustomerLocation) {
	location := CustomerLocation{
		ID:              row.Int(res.Map("locationID")),
		CustomerID:      row.Int(res.Map("cust_id")),
		Name:            row.Str(res.Map("name")),
		Address:         row.Str(res.Map("address")),
		City:            row.Str(res.Map("city")),
		StateID:         row.Int(res.Map("stateID")),
		PostalCode:      row.Str(res.Map("postalCode")),
		Email:           row.Str(res.Map("email")),
		Phone:           row.Str(res.Map("phone")),
		Fax:             row.Str(res.Map("fax")),
		Latitude:        row.Float(res.Map("latitude")),
		Longitude:       row.Float(res.Map("longitude")),
		ContactPerson:   row.Str(res.Map("contact_person")),
		IsPrimary:       row.Bool(res.Map("isprimary")),
		ShippingDefault: row.Bool(res.Map("ShippingDefault")),
	}
	ch <- location
}
