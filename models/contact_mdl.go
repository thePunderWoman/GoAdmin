package models

import (
	"../helpers/database"
	_ "errors"
	"github.com/ziutek/mymysql/mysql"
	_ "log"
	"time"
)

type Contact struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	Subject     string
	Message     string
	Created     time.Time
	Type        string
	ContactType ContactType
	Address1    string
	Address2    string
	City        string
	State       string
	PostalCode  string
	Country     string
}

type ContactType struct {
	ID        int
	Name      string
	Receivers ContactReceivers
}

type ContactReceiver struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
}

type ContactReceivers []ContactReceiver

func (c Contact) GetAll() ([]Contact, error) {
	contacts := make([]Contact, 0)
	sel, err := database.GetStatement("getAllContactsStmt")
	if err != nil {
		return contacts, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return contacts, err
	}

	ch := make(chan Contact)
	for _, row := range rows {
		go c.PopulateContact(row, res, ch)
	}
	for _, _ = range rows {
		contacts = append(contacts, <-ch)
	}
	return contacts, nil
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

func (c *Contact) Get() error {
	sel, err := database.GetStatement("getContactStmt")
	if err != nil {
		return err
	}
	sel.Bind(c.ID)
	row, res, err := sel.ExecFirst()
	ch := make(chan Contact)
	go c.PopulateContact(row, res, ch)
	contact := <-ch
	c.FirstName = contact.FirstName
	c.LastName = contact.LastName
	c.Email = contact.Email
	c.Type = contact.Type
	c.Phone = contact.Phone
	c.Subject = contact.Subject
	c.Message = contact.Message
	c.Created = contact.Created
	c.Address1 = contact.Address1
	c.Address2 = contact.Address2
	c.City = contact.City
	c.PostalCode = contact.PostalCode
	c.Country = contact.Country
	return nil
}

func (t ContactType) GetAll() ([]ContactType, error) {
	types := make([]ContactType, 0)
	sel, err := database.GetStatement("getAllContactTypesStmt")
	if err != nil {
		return types, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return types, err
	}

	ch := make(chan ContactType)
	for _, row := range rows {
		go t.PopulateContactType(row, res, ch)
	}
	for _, _ = range rows {
		types = append(types, <-ch)
	}

	return types, nil
}

func (t ContactType) PopulateContactType(row mysql.Row, res mysql.Result, ch chan ContactType) {
	ctype := ContactType{
		ID:   row.Int(res.Map("contactTypeID")),
		Name: row.Str(res.Map("name")),
	}
	ch <- ctype
}
