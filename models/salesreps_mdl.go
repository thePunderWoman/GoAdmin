package models

import (
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	"log"
)

type SalesRep struct {
	ID            int
	Name          string
	Code          string
	CustomerCount int
}

func (s SalesRep) GetAll() ([]SalesRep, error) {
	reps := make([]SalesRep, 0)
	sel, err := database.GetStatement("GetAllSalesRepsStmt")
	if err != nil {
		return reps, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return reps, err
	}
	ch := make(chan SalesRep)
	for _, row := range rows {
		go s.PopulateSalesRep(row, res, ch)
	}

	for _, _ = range rows {
		reps = append(reps, <-ch)
	}

	return reps, nil
}

func (s SalesRep) Get() (SalesRep, error) {
	var rep SalesRep
	sel, err := database.GetStatement("GetSalesRepStmt")
	if err != nil {
		return rep, err
	}
	sel.Bind(s.ID)
	row, res, err := sel.ExecFirst()
	if err != nil {
		return rep, err
	}
	ch := make(chan SalesRep)
	go s.PopulateSalesRep(row, res, ch)
	return <-ch, nil
}

func (s SalesRep) PopulateSalesRep(row mysql.Row, res mysql.Result, ch chan SalesRep) {
	rep := SalesRep{
		ID:            row.Int(res.Map("salesRepID")),
		Name:          row.Str(res.Map("name")),
		Code:          row.Str(res.Map("code")),
		CustomerCount: row.Int(res.Map("customercount")),
	}
	ch <- rep
}

func (s *SalesRep) Save() error {
	if s.ID > 0 {
		// update
		upd, err := database.GetStatement("UpdateSalesRepStmt")
		if err != nil {
			return err
		}
		upd.Bind(s.Name, s.Code, s.ID)
		_, _, err = upd.Exec()
		return err
	} else {
		// new
		ins, err := database.GetStatement("AddSalesRepStmt")
		if err != nil {
			return err
		}
		ins.Bind(s.Name, s.Code)
		_, _, err = ins.Exec()
		return err
	}
	return nil
}

func (s SalesRep) Delete() bool {
	del, err := database.GetStatement("DeleteSalesRepStmt")
	if err != nil {
		log.Println(err)
		return false
	}
	del.Bind(s.ID)
	_, _, err = del.Exec()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
