package models

import (
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	_ "log"
	"time"
)

var CentralTime, _ = time.LoadLocation("US/Central")

type NewsItem struct {
	ID           int
	Title        string
	Lead         string
	Content      string
	PublishStart time.Time
	PublishEnd   time.Time
	Active       bool
	Slug         string
}

type News []NewsItem

func (n NewsItem) GetAll() (News, error) {
	var news News
	sel, err := database.GetStatement("GetAllNewsStmt")
	if err != nil {
		return news, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return news, err
	}
	ch := make(chan NewsItem)
	for _, row := range rows {
		go n.PopulateNewsItem(row, res, ch)
	}

	for _, _ = range rows {
		news = append(news, <-ch)
	}

	return news, nil
}

func (n NewsItem) PopulateNewsItem(row mysql.Row, res mysql.Result, ch chan NewsItem) {
	item := NewsItem{
		ID:           row.Int(res.Map("newsItemID")),
		Title:        row.Str(res.Map("title")),
		Lead:         row.Str(res.Map("lead")),
		Content:      row.Str(res.Map("content")),
		PublishStart: row.Time(res.Map("publishStart"), UTC),
		PublishEnd:   row.Time(res.Map("publishEnd"), UTC),
		Active:       row.Bool(res.Map("active")),
		Slug:         row.Str(res.Map("slug")),
	}
	ch <- item
}

func (n NewsItem) Get() (NewsItem, error) {
	var newsitem NewsItem
	sel, err := database.GetStatement("GetNewsItemStmt")
	if err != nil {
		return newsitem, err
	}
	sel.Bind(n.ID)
	row, res, err := sel.ExecFirst()
	if err != nil {
		return newsitem, err
	}
	ch := make(chan NewsItem)
	go n.PopulateNewsItem(row, res, ch)
	newsitem = <-ch
	return newsitem, nil
}

func (n *NewsItem) Save() error {
	n.SetZones()
	if n.ID > 0 {
		// update
		upd, err := database.GetStatement("UpdateNewsItemStmt")
		if err != nil {
			return err
		}

		upd.Bind(n.Title, n.Lead, n.Content, n.PublishStart, n.PublishEnd, GenerateSlug(n.Title), n.ID)
		_, _, err = upd.Exec()
		return err
	} else {
		// new
		ins, err := database.GetStatement("AddNewsItemStmt")
		if err != nil {
			return err
		}
		ins.Bind(n.Title, n.Lead, n.Content, n.PublishStart, n.PublishEnd, GenerateSlug(n.Title))
		_, _, err = ins.Exec()
		return err
	}
	return nil
}

func (n NewsItem) Delete() error {
	del, err := database.GetStatement("DeleteNewsItemStmt")
	if err != nil {
		return err
	}
	del.Bind(n.ID)
	_, _, err = del.Exec()
	return err
}

func (n *NewsItem) SetZones() {
	if !n.PublishStart.IsZero() {
		n.PublishStart = ChangeZone(n.PublishStart, CentralTime)
	}
	if !n.PublishEnd.IsZero() {
		n.PublishEnd = ChangeZone(n.PublishEnd, CentralTime)
	}
}

func ChangeZone(t time.Time, zone *time.Location) time.Time {
	if !t.IsZero() {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), zone).In(UTC)
	}
	return t
}
