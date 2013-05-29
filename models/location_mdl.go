package models

import (
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	//"log"
	"sort"
	//"strconv"
	//"time"
)

type State struct {
	ID        int
	Name      string
	Abbr      string
	CountryID int
}

type States []State

type Country struct {
	ID     int
	Name   string
	Abbr   string
	States States
}

type Countries []Country

func (c Country) GetAll() (Countries, error) {
	var countries Countries
	sel, err := database.GetStatement("GetAllCountriesStmt")
	if err != nil {
		return countries, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return countries, err
	}
	ch := make(chan Country)
	for _, row := range rows {
		go c.PopulateCountry(row, res, ch)
	}
	for _, _ = range rows {
		country := <-ch
		country.GetStates()
		countries = append(countries, country)
	}
	countries.Sort()
	return countries, nil
}

func (c Country) PopulateCountry(row mysql.Row, res mysql.Result, ch chan Country) {
	country := Country{
		ID:   row.Int(res.Map("countryID")),
		Name: row.Str(res.Map("name")),
		Abbr: row.Str(res.Map("abbr")),
	}
	ch <- country
}

func (c *Country) GetStates() error {
	sel, err := database.GetStatement("GetStatesByCountryStmt")
	if err != nil {
		return err
	}
	sel.Bind(c.ID)
	rows, res, err := sel.Exec()
	if err != nil {
		return err
	}
	var states States
	ch := make(chan State)
	s := State{}
	for _, row := range rows {
		go s.PopulateState(row, res, ch)
	}
	for _, _ = range rows {
		states = append(states, <-ch)
	}
	states.Sort()
	c.States = states
	return nil
}

func (s State) GetAll() (States, error) {
	var states States
	sel, err := database.GetStatement("GetAllStatesStmt")
	if err != nil {
		return states, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return states, err
	}
	ch := make(chan State)
	for _, row := range rows {
		go s.PopulateState(row, res, ch)
	}
	for _, _ = range rows {
		states = append(states, <-ch)
	}
	return states, nil
}

func (s State) Get() (State, error) {
	var state State
	if s.ID > 0 {
		sel, err := database.GetStatement("GetStateStmt")
		if err != nil {
			return state, err
		}
		sel.Raw.Reset()
		sel.Bind(s.ID)
		row, res, err := sel.ExecFirst()
		if err != nil {
			return state, err
		}
		ch := make(chan State)
		go s.PopulateState(row, res, ch)
		state = <-ch
	}
	return state, nil
}

func (s State) PopulateState(row mysql.Row, res mysql.Result, ch chan State) {
	state := State{
		ID:        row.Int(res.Map("stateID")),
		CountryID: row.Int(res.Map("countryID")),
		Name:      row.Str(res.Map("state")),
		Abbr:      row.Str(res.Map("abbr")),
	}
	ch <- state
}

func (c Countries) Len() int           { return len(c) }
func (c Countries) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Countries) Less(i, j int) bool { return c[i].Name < (c[j].Name) }

func (c *Countries) Sort() {
	sort.Sort(c)
}

func (s States) Len() int           { return len(s) }
func (s States) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s States) Less(i, j int) bool { return s[i].Name < (s[j].Name) }

func (s *States) Sort() {
	sort.Sort(s)
}

func (s States) ToMap() map[int]State {
	statemap := make(map[int]State, 0)
	for _, state := range s {
		statemap[state.ID] = state
	}
	return statemap
}
