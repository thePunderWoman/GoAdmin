package models

import (
	//"../helpers/UDF"
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	"time"
	//"sort"
)

type CustomerUser struct {
	ID                 string
	CustID, CustomerID int
	Name               string
	Email              string
	Password           string
	DateAdded          time.Time
	Active             bool
	LocationID         int
	IsSudo             bool
	NotCustomer        bool
	Keys               APIKeys
}

type APIKey struct {
	ID        string
	UserID    string
	Key       string
	TypeID    string
	DateAdded time.Time
	KeyType   APIKeyType
}

type APIKeys []APIKey

type APIKeyType struct {
	ID, Type  string
	DateAdded time.Time
}

func (c CustomerUser) GetAll() (users []CustomerUser, err error) {
	keyMap := make(map[string]APIKeys)
	keychan := make(chan int)

	go func(ch chan int) {
		keys, _ := c.GetCustomerKeys()
		keyMap = keys.ToMap()
		ch <- 1
	}(keychan)

	sel, err := database.GetStatement("GetAllCustomerUsersStmt")
	if err != nil {
		return users, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return users, err
	}
	<-keychan

	ch := make(chan CustomerUser)
	for _, row := range rows {
		go c.PopulateUser(row, res, ch)
	}
	for _, _ = range rows {
		custuser := <-ch
		custuser.Keys = keyMap[custuser.ID]
		users = append(users, custuser)
	}
	return
}

func (c CustomerUser) GetAllByCustomer() (users []CustomerUser, err error) {

	keyMap := make(map[string]APIKeys)
	keychan := make(chan int)

	go func(ch chan int) {
		keys, _ := c.GetCustomerKeys()
		keyMap = keys.ToMap()
		ch <- 1
	}(keychan)

	sel, err := database.GetStatement("GetCustomerUsersStmt")
	if err != nil {
		return users, err
	}
	sel.Bind(c.CustID)
	rows, res, err := sel.Exec()
	if err != nil {
		return users, err
	}
	<-keychan

	ch := make(chan CustomerUser)
	for _, row := range rows {
		go c.PopulateUser(row, res, ch)
	}
	for _, _ = range rows {
		custuser := <-ch
		custuser.Keys = keyMap[custuser.ID]
		users = append(users, custuser)
	}
	return
}

func (c CustomerUser) PopulateUser(row mysql.Row, res mysql.Result, ch chan CustomerUser) {
	user := CustomerUser{
		ID:          row.Str(res.Map("id")),
		CustID:      row.Int(res.Map("cust_ID")),
		CustomerID:  row.Int(res.Map("customerID")),
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

func (c CustomerUser) GetAllCustomerKeys() (keys APIKeys, err error) {
	sel, err := database.GetStatement("GetAllCustomerUserKeysStmt")
	if err != nil {
		return keys, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return keys, err
	}

	ch := make(chan APIKey)
	for _, row := range rows {
		go c.PopulateKey(row, res, ch)
	}
	for _, _ = range rows {
		keys = append(keys, <-ch)
	}
	return
}

func (c CustomerUser) GetCustomerKeys() (keys APIKeys, err error) {
	sel, err := database.GetStatement("GetCustomerUserKeysStmt")
	if err != nil {
		return keys, err
	}
	sel.Bind(c.CustID)
	rows, res, err := sel.Exec()
	if err != nil {
		return keys, err
	}

	ch := make(chan APIKey)
	for _, row := range rows {
		go c.PopulateKey(row, res, ch)
	}
	for _, _ = range rows {
		keys = append(keys, <-ch)
	}
	return
}

func (c CustomerUser) PopulateKey(row mysql.Row, res mysql.Result, ch chan APIKey) {
	keyType := APIKeyType{
		ID:        row.Str(res.Map("type_id")),
		Type:      row.Str(res.Map("type")),
		DateAdded: row.Time(res.Map("typeDateAdded"), UTC),
	}
	key := APIKey{
		ID:        row.Str(res.Map("id")),
		UserID:    row.Str(res.Map("user_id")),
		Key:       row.Str(res.Map("api_key")),
		TypeID:    keyType.ID,
		DateAdded: row.Time(res.Map("date_added"), UTC),
		KeyType:   keyType,
	}
	ch <- key
}

func (k APIKeys) ToMap() map[string]APIKeys {
	keymap := make(map[string]APIKeys, 0)
	for _, key := range k {
		keymap[key.UserID] = append(keymap[key.UserID], key)
	}
	return keymap
}
