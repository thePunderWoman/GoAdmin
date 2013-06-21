package models

import (
	//"../helpers/UDF"
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	"time"
	//"sort"
)

type CustomerUser struct {
	ID          int
	CustID      int
	Name        string
	Email       string
	Password    string
	DateAdded   time.Time
	Active      bool
	LocationID  int
	IsSudo      bool
	NotCustomer bool
	Keys        []APIkey
}

type APIKey struct {
	ID        int
	Key       string
	TypeID    int
	DateAdded time.Time
}

func (c CustomerUser) GetAllByCustomer() (users []CustomerUser, err error) {

	sel, err := database.GetStatement("GetCustomerUsersStmt")
	if err != nil {
		return users, err
	}
	sel.Bind(c.CustID)
	rows, res, err := sel.Exec()
	if err != nil {
		return users, err
	}

	ch := make(chan CustomerUser)
	for _, row := range rows {
		id := row.Int(res.Map("id"))

		go c.PopulateUser(row, res, ch)
	}
	for _, _ = range rows {
		users = append(users, <-ch)
	}
	return
}

func (c CustomerUser) PopulateUser(row mysql.Row, res mysql.Result, ch chan CustomerUser) {
	user := CustomerUser{
		ID:          row.Int(res.Map("id")),
		CustID:      row.Int(res.Map("cust_ID")),
		Name:        row.Str(res.Map("name")),
		Email:       row.Str(res.Map("email")),
		DateAdded:   row.Time(res.Map("date_added"), UTC),
		Active:      row.Bool(res.Map("active")),
		LocationID:  row.Int(res.Map("locationID")),
		IsSudo:      row.Bool(res.Map("isSudo")),
		NotCustomer: row.Bool(res.Map("NotCustomer")),
	}
	ch <- user
}
