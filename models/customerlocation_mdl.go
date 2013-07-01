package models

import (
	"../helpers/UDF"
	"../helpers/database"
	"../helpers/geo"
	"github.com/ziutek/mymysql/mysql"
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

func (c CustomerLocation) GetAllNoGeo() (locations CustomerLocations, err error) {
	stateMap := make(map[int]State)
	statechan := make(chan int)

	go func(ch chan int) {
		states, _ := State{}.GetAll()
		stateMap = states.ToMap()
		ch <- 1
	}(statechan)

	sel, err := database.GetStatement("GetCustomerLocationsNoGeoStmt")
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

func (c CustomerLocation) Get() (location CustomerLocation, err error) {
	sel, err := database.GetStatement("GetCustomerLocationStmt")
	if err != nil {
		return location, err
	}
	sel.Bind(c.ID)
	row, res, err := sel.ExecFirst()
	if err != nil {
		return location, err
	}

	ch := make(chan CustomerLocation)
	go c.PopulateLocation(row, res, ch)
	location = <-ch
	location.State, _ = State{ID: location.StateID}.Get()
	return
}

func (c CustomerLocation) Save() error {
	if c.ID > 0 {
		// update location
		upd, err := database.GetStatement("UpdateCustomerLocationStmt")
		if err != nil {
			return err
		}
		if c.Latitude == 0 && c.Longitude == 0 {
			state, _ := State{ID: c.StateID}.Get()
			latlng, err := geo.GetLatLng(c.Address, c.City, state.Abbr)
			if err == nil {
				c.Latitude = latlng.Lat
				c.Longitude = latlng.Lng
			}
		}
		params := struct {
			Name       string
			Address    string
			City       string
			StateID    int
			PostalCode string
			Email      *string
			Phone      *string
			Fax        *string
			Latitude   float64
			Longitude  float64
			ID         int
		}{
			c.Name,
			c.Address,
			c.City,
			c.StateID,
			c.PostalCode,
			UDF.StrOrNil(c.Email),
			UDF.StrOrNil(c.Phone),
			UDF.StrOrNil(c.Fax),
			c.Latitude,
			c.Longitude,
			c.ID,
		}
		upd.Bind(&params)
		_, _, err = upd.Exec()
		return err

	} else {
		// new location
		ins, err := database.GetStatement("AddCustomerLocationStmt")
		if err != nil {
			return err
		}
		if c.Latitude == 0 && c.Longitude == 0 {
			state, _ := State{ID: c.StateID}.Get()
			latlng, err := geo.GetLatLng(c.Address, c.City, state.Abbr)
			if err == nil {
				c.Latitude = latlng.Lat
				c.Longitude = latlng.Lng
			}
		}
		params := struct {
			Name            string
			Address         string
			City            string
			StateID         int
			PostalCode      string
			Email           *string
			Phone           *string
			Fax             *string
			Latitude        float64
			Longitude       float64
			CustomerID      int
			IsPrimary       bool
			ShippingDefault bool
			ContactPerson   *string
		}{
			c.Name,
			c.Address,
			c.City,
			c.StateID,
			c.PostalCode,
			UDF.StrOrNil(c.Email),
			UDF.StrOrNil(c.Phone),
			UDF.StrOrNil(c.Fax),
			c.Latitude,
			c.Longitude,
			c.CustomerID,
			false,
			false,
			nil,
		}
		ins.Bind(&params)
		_, _, err = ins.Exec()
		return err

	}
	return nil
}

func (c CustomerLocation) Delete() bool {
	del, err := database.GetStatement("DeleteCustomerLocationStmt")
	if err != nil {
		return false
	}
	del.Bind(c.ID)
	_, _, err = del.Exec()
	if err != nil {
		return false
	}
	return true
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
