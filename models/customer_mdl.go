package models

import (
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	//"log"
	"sort"
	//"strconv"
	//"time"
)

type Customer struct {
	ID            int
	Name          string
	Email         string
	Address       string
	Address2      string
	City          string
	StateID       int
	State         State
	PostalCode    string
	Phone         string
	Fax           string
	ContactPerson string
	DealerTypeID  int
	DealerType    DealerType
	Latitude      string
	Longitude     string
	Website       string
	CustomerID    int
	IsDummy       bool
	ParentID      int
	SearchURL     string
	ELocalURL     string
	Logo          string
	MapicsCodeID  int
	MapicsCode    MapicsCode
	SalesRepID    int
	SalesRep      SalesRep
	APIKey        string
	Tier          int
	DealerTier    DealerTier
	ShowWebsite   bool
}

type DealerType struct {
	ID     int
	Name   string
	Online bool
	Show   bool
	Label  string
}
type DealerTypes []DealerType

type DealerTier struct {
	ID   int
	Name string
	Sort int
}
type DealerTiers []DealerTier

type MapicsCode struct {
	ID          int
	Code        string
	Description string
}
type MapicsCodes []MapicsCode

type MapIcons struct {
	ID            int
	DealerTierID  int
	DealerTypeID  int
	MapIcon       string
	MapIconShadow string
}

type Customers []Customer

func (c Customer) GetAll() (Customers, error) {
	var customers Customers
	sel, err := database.GetStatement("GetAllCustomersStmt")
	if err != nil {
		return customers, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return customers, err
	}

	ch := make(chan Customer)
	for _, row := range rows {
		go c.PopulateCustomer(row, res, ch)
	}
	for _, _ = range rows {
		customers = append(customers, <-ch)
	}
	return customers, nil
}

func (c Customer) Get() (cust Customer, err error) {
	sel, err := database.GetStatement("GetCustomerStmt")
	if err != nil {
		return
	}
	sel.Bind(c.ID)
	row, res, err := sel.ExecFirst()
	if err != nil {
		return
	}

	ch := make(chan Customer)
	go c.PopulateCustomer(row, res, ch)
	cust = <-ch
	return
}

func (c Customer) PopulateCustomer(row mysql.Row, res mysql.Result, ch chan Customer) {
	customer := Customer{
		ID:            row.Int(res.Map("cust_id")),
		Name:          row.Str(res.Map("name")),
		Email:         row.Str(res.Map("email")),
		Address:       row.Str(res.Map("address")),
		Address2:      row.Str(res.Map("address2")),
		City:          row.Str(res.Map("city")),
		StateID:       row.Int(res.Map("stateID")),
		PostalCode:    row.Str(res.Map("postal_code")),
		Phone:         row.Str(res.Map("phone")),
		Fax:           row.Str(res.Map("fax")),
		ContactPerson: row.Str(res.Map("contact_person")),
		DealerTypeID:  row.Int(res.Map("dealer_type")),
		Latitude:      row.Str(res.Map("latitude")),
		Longitude:     row.Str(res.Map("longitude")),
		Website:       row.Str(res.Map("website")),
		CustomerID:    row.Int(res.Map("customerID")),
		IsDummy:       row.Bool(res.Map("isDummy")),
		ParentID:      row.Int(res.Map("parentID")),
		SearchURL:     row.Str(res.Map("searchURL")),
		ELocalURL:     row.Str(res.Map("eLocalURL")),
		Logo:          row.Str(res.Map("logo")),
		MapicsCodeID:  row.Int(res.Map("mCodeID")),
		SalesRepID:    row.Int(res.Map("salesRepID")),
		APIKey:        row.Str(res.Map("APIKey")),
		Tier:          row.Int(res.Map("tier")),
		ShowWebsite:   row.Bool(res.Map("showWebsite")),
	}
	statechan := make(chan State)
	typechan := make(chan DealerType)
	tierchan := make(chan DealerTier)
	codechan := make(chan MapicsCode)
	repchan := make(chan SalesRep)
	go func(id int, ch chan State) {
		state := State{ID: id}
		state, _ = state.Get()
		ch <- state
	}(customer.StateID, statechan)
	go func(id int, ch chan DealerType) {
		dealertype := DealerType{ID: id}
		dealertype, _ = dealertype.Get()
		ch <- dealertype
	}(customer.DealerTypeID, typechan)
	go func(id int, ch chan DealerTier) {
		tier := DealerTier{ID: id}
		tier, _ = tier.Get()
		ch <- tier
	}(customer.Tier, tierchan)
	go func(id int, ch chan MapicsCode) {
		code := MapicsCode{ID: id}
		code, _ = code.Get()
		ch <- code
	}(customer.MapicsCodeID, codechan)
	go func(id int, ch chan SalesRep) {
		rep := SalesRep{ID: id}
		rep, _ = rep.Get()
		ch <- rep
	}(customer.SalesRepID, repchan)

	customer.State = <-statechan
	customer.DealerType = <-typechan
	customer.DealerTier = <-tierchan
	customer.MapicsCode = <-codechan
	customer.SalesRep = <-repchan

	ch <- customer
}

func (d DealerType) GetAll() (types DealerTypes, err error) {
	sel, err := database.GetStatement("GetAllDealerTypesStmt")
	if err != nil {
		return
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return
	}
	ch := make(chan DealerType)
	for _, row := range rows {
		go d.PopulateType(row, res, ch)
	}
	for _, _ = range rows {
		types = append(types, <-ch)
	}
	return
}

func (d DealerType) Get() (dealertype DealerType, err error) {
	sel, err := database.GetStatement("GetDealerTypeStmt")
	if err != nil {
		return
	}
	sel.Bind(d.ID)
	row, res, err := sel.ExecFirst()
	if err != nil {
		return
	}
	ch := make(chan DealerType)
	go d.PopulateType(row, res, ch)
	dealertype = <-ch
	return
}

func (d DealerType) PopulateType(row mysql.Row, res mysql.Result, ch chan DealerType) {
	dealertype := DealerType{
		ID:     row.Int(res.Map("dealer_type")),
		Name:   row.Str(res.Map("type")),
		Online: row.Bool(res.Map("online")),
		Show:   row.Bool(res.Map("show")),
		Label:  row.Str(res.Map("label")),
	}
	ch <- dealertype
}

func (d DealerTier) GetAll() (tiers DealerTiers, err error) {
	sel, err := database.GetStatement("GetAllDealerTiersStmt")
	if err != nil {
		return
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return
	}
	ch := make(chan DealerTier)
	for _, row := range rows {
		go d.PopulateTier(row, res, ch)
	}
	for _, _ = range rows {
		tiers = append(tiers, <-ch)
	}
	return
}

func (d DealerTier) Get() (tier DealerTier, err error) {
	sel, err := database.GetStatement("GetDealerTierStmt")
	if err != nil {
		return
	}
	sel.Bind(d.ID)
	row, res, err := sel.ExecFirst()
	if err != nil {
		return
	}
	ch := make(chan DealerTier)
	go d.PopulateTier(row, res, ch)
	tier = <-ch
	return
}

func (d DealerTier) PopulateTier(row mysql.Row, res mysql.Result, ch chan DealerTier) {
	tier := DealerTier{
		ID:   row.Int(res.Map("ID")),
		Name: row.Str(res.Map("tier")),
		Sort: row.Int(res.Map("sort")),
	}
	ch <- tier
}

func (m MapicsCode) GetAll() (codes MapicsCodes, err error) {
	sel, err := database.GetStatement("GetAllMapicsCodesStmt")
	if err != nil {
		return
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return
	}
	ch := make(chan MapicsCode)
	for _, row := range rows {
		go m.PopulateCode(row, res, ch)
	}
	for _, _ = range rows {
		codes = append(codes, <-ch)
	}
	return
}

func (m MapicsCode) Get() (code MapicsCode, err error) {
	if m.ID > 0 {
		sel, err := database.GetStatement("GetMapicsCodeStmt")
		if err != nil {
			return code, err
		}
		sel.Raw.Reset()
		sel.Bind(m.ID)
		row, res, err := sel.ExecFirst()
		if err != nil {
			return code, err
		}
		ch := make(chan MapicsCode)
		go m.PopulateCode(row, res, ch)
		code = <-ch
	}
	return
}

func (m MapicsCode) PopulateCode(row mysql.Row, res mysql.Result, ch chan MapicsCode) {
	code := MapicsCode{
		ID:          row.Int(res.Map("mCodeID")),
		Code:        row.Str(res.Map("code")),
		Description: row.Str(res.Map("description")),
	}
	ch <- code
}

func (c Customers) Len() int           { return len(c) }
func (c Customers) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Customers) Less(i, j int) bool { return c[i].Name < (c[j].Name) }

func (c *Customers) Sort() {
	sort.Sort(c)
}

func (d DealerTypes) Len() int           { return len(d) }
func (d DealerTypes) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DealerTypes) Less(i, j int) bool { return d[i].Name < (d[j].Name) }

func (d *DealerTypes) Sort() {
	sort.Sort(d)
}

func (d DealerTiers) Len() int           { return len(d) }
func (d DealerTiers) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DealerTiers) Less(i, j int) bool { return d[i].Name < (d[j].Name) }

func (d *DealerTiers) Sort() {
	sort.Sort(d)
}

func (m MapicsCodes) Len() int           { return len(m) }
func (m MapicsCodes) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MapicsCodes) Less(i, j int) bool { return m[i].Code < (m[j].Code) }

func (m *MapicsCodes) Sort() {
	sort.Sort(m)
}
