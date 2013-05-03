package models

import (
	"../helpers/database"
	_ "errors"
	"github.com/ziutek/mymysql/mysql"
	_ "log"
	"sort"
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
	ID   int
	Name string
}

type ContactReceiver struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Types     ContactTypes
}

type ContactTypes []ContactType

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
		go PopulateContactType(row, res, ch)
	}
	for _, _ = range rows {
		types = append(types, <-ch)
	}

	return types, nil
}

func PopulateContactType(row mysql.Row, res mysql.Result, ch chan ContactType) {
	ctype := ContactType{
		ID:   row.Int(res.Map("contactTypeID")),
		Name: row.Str(res.Map("name")),
	}
	ch <- ctype
}

func (r ContactReceiver) GetAll() ([]ContactReceiver, error) {
	receivers := make([]ContactReceiver, 0)
	sel, err := database.GetStatement("getAllContactReceiversStmt")
	if err != nil {
		return receivers, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return receivers, err
	}

	ch := make(chan ContactReceiver)
	for _, row := range rows {
		go r.PopulateContactReceiver(row, res, ch)
	}
	for _, _ = range rows {
		receivers = append(receivers, <-ch)
	}

	return receivers, nil
}

func (r ContactReceiver) PopulateContactReceiver(row mysql.Row, res mysql.Result, ch chan ContactReceiver) {
	receiver := ContactReceiver{
		ID:        row.Int(res.Map("contactReceiverID")),
		FirstName: row.Str(res.Map("first_name")),
		LastName:  row.Str(res.Map("last_name")),
		Email:     row.Str(res.Map("email")),
	}
	receiver.GetTypes()
	ch <- receiver
}

func (r *ContactReceiver) GetTypes() {
	var types ContactTypes
	sel, err := database.GetStatement("getReceiverContactTypesStmt")
	if err != nil {
		r.Types = types
		return
	}
	sel.Reset()
	sel.Bind(r.ID)
	rows, res, err := sel.Exec()
	if err != nil {
		r.Types = types
		return
	}
	ch := make(chan ContactType)
	for _, row := range rows {
		go PopulateContactType(row, res, ch)
	}
	for _, _ = range rows {
		types = append(types, <-ch)
	}
	types.Sort()
	r.Types = types
}

func (t ContactTypes) Len() int           { return len(t) }
func (t ContactTypes) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ContactTypes) Less(i, j int) bool { return t[i].Name < (t[j].Name) }

func (t *ContactTypes) Sort() {
	sort.Sort(t)
}
