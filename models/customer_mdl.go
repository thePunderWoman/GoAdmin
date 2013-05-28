package models

import (
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	_ "log"
	"sort"
	"strconv"
	"time"
)

type Customer struct {
	ID            int
	Name          string
	Email         string
	Address       string
	Address2      string
	City          string
	StateID       string
	PostalCode    string
	Phone         string
	Fax           string
	ContactPerson string
	DealerTypeID  int
	Latitude      string
	Longitude     string
	Website       string
	CustomerID    string
	IsDummy       bool
	ParentID      int
	SearchURL     string
	ELocalURL     string
	Logo          string
	MapicsCodeID  int
	SalesRepID    int
	APIKey        string
	Tier          int
	ShowWebsite   bool
}

type Customers []Customer

func (c Customer) GetAll() (Customers, error) {
	var customers Customers
	sel, err := database.GetStatement("getAllContactsStmt")
	if err != nil {
		return customers, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return customers, err
	}

	ch := make(chan Contact)
	for _, row := range rows {
		go c.PopulateContact(row, res, ch)
	}
	for _, _ = range rows {
		customers = append(customers, <-ch)
	}
	return customers, nil
}

func (c Contact) PopulateContact(row mysql.Row, res mysql.Result, ch chan Contact) {
	contact := Contact{
		ID:         row.Int(res.Map("contactID")),
		FirstName:  row.Str(res.Map("first_name")),
		LastName:   row.Str(res.Map("last_name")),
		Email:      row.Str(res.Map("email")),
		Phone:      row.Str(res.Map("phone")),
		Subject:    row.Str(res.Map("subject")),
		Message:    row.Str(res.Map("message")),
		Created:    row.Time(res.Map("createdDate"), UTC),
		Type:       row.Str(res.Map("type")),
		Address1:   row.Str(res.Map("address1")),
		Address2:   row.Str(res.Map("address2")),
		City:       row.Str(res.Map("city")),
		PostalCode: row.Str(res.Map("postalcode")),
		Country:    row.Str(res.Map("country")),
	}
	ch <- contact
}

func (c Customers) Len() int           { return len(c) }
func (c Customers) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Customers) Less(i, j int) bool { return c[i].Name < (c[j].Name) }

func (c *Customers) Sort() {
	sort.Sort(c)
}
