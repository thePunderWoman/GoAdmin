package models

import (
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	_ "log"
	"sort"
)

type FAQ struct {
	ID       int
	Question string
	Answer   string
}
type FAQS []FAQ

func (f FAQ) GetAll() (FAQS, error) {
	faqs := FAQS{}
	sel, err := database.GetStatement("GetAllFAQStmt")
	if err != nil {
		return faqs, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return faqs, err
	}
	ch := make(chan FAQ)
	for _, row := range rows {
		go f.PopulateFAQ(row, res, ch)
	}

	for _, _ = range rows {
		faqs = append(faqs, <-ch)
	}
	faqs.Sort()
	return faqs, nil
}

func (f FAQ) Get() (FAQ, error) {
	faq := FAQ{}
	sel, err := database.GetStatement("GetFAQStmt")
	if err != nil {
		return faq, err
	}
	sel.Bind(f.ID)
	row, res, err := sel.ExecFirst()
	if err != nil {
		return faq, err
	}
	ch := make(chan FAQ)
	go f.PopulateFAQ(row, res, ch)
	faq = <-ch

	return faq, nil
}

func (f *FAQ) Save() error {
	if f.ID > 0 {
		//update
		upd, err := database.GetStatement("UpdateFAQStmt")
		if err != nil {
			return err
		}
		upd.Bind(f.Question, f.Answer, f.ID)
		_, _, err = upd.Exec()
		if err != nil {
			return err
		}
	} else {
		// new
		ins, err := database.GetStatement("AddFAQStmt")
		if err != nil {
			return err
		}
		ins.Bind(f.Question, f.Answer)
		_, _, err = ins.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FAQ) Delete() error {
	del, err := database.GetStatement("DeleteFAQStmt")
	if err != nil {
		return err
	}
	del.Bind(f.ID)
	_, _, err = del.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (f *FAQ) PopulateFAQ(row mysql.Row, res mysql.Result, ch chan FAQ) {
	faq := FAQ{
		ID:       row.Int(res.Map("faqID")),
		Question: row.Str(res.Map("question")),
		Answer:   row.Str(res.Map("answer")),
	}
	ch <- faq
}

func (f FAQS) Len() int           { return len(f) }
func (f FAQS) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FAQS) Less(i, j int) bool { return f[i].Question < (f[j].Question) }

func (f *FAQS) Sort() {
	sort.Sort(f)
}
