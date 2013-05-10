package models

import (
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	"log"
	"sort"
	_ "strconv"
	"time"
)

type LandingPage struct {
	ID              int
	Name            string
	Start           time.Time
	End             time.Time
	URL             string
	Content         string
	LinkClasses     string
	ConversionID    string
	ConversionLabel string
	NewWindow       bool
	MenuPosition    string
	Images          LandingPageImages
	Data            []LandingPageData
}

type LandingPageImage struct {
	ID            int
	LandingPageID int
	URL           string
	Sort          int
}

type LandingPageImages []LandingPageImage

type LandingPageData struct {
	ID            int
	LandingPageID int
	Key           string
	Value         string
}

func (l LandingPageImages) Len() int           { return len(l) }
func (l LandingPageImages) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l LandingPageImages) Less(i, j int) bool { return l[i].Sort < (l[j].Sort) }

func (l *LandingPageImages) Sort() {
	sort.Sort(l)
}

func (l LandingPage) GetActive() ([]LandingPage, error) {
	pages := make([]LandingPage, 0)
	sel, err := database.GetStatement("GetActiveLandingPagesStmt")
	if err != nil {
		return pages, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return pages, err
	}
	ch := make(chan LandingPage)
	for _, row := range rows {
		go l.populateLandingPage(row, res, ch)
	}
	for _, _ = range rows {
		pages = append(pages, <-ch)
	}
	return pages, nil
}

func (l LandingPage) GetPast() ([]LandingPage, error) {
	pages := make([]LandingPage, 0)
	sel, err := database.GetStatement("GetPastLandingPagesStmt")
	if err != nil {
		return pages, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return pages, err
	}
	ch := make(chan LandingPage)
	for _, row := range rows {
		go l.populateLandingPage(row, res, ch)
	}
	for _, _ = range rows {
		pages = append(pages, <-ch)
	}
	return pages, nil
}

func (l LandingPage) Get() (LandingPage, error) {
	var page LandingPage
	sel, err := database.GetStatement("GetLandingPageStmt")
	if err != nil {
		return page, err
	}
	sel.Bind(l.ID)
	row, res, err := sel.ExecFirst()
	if err != nil {
		return page, err
	}
	ch := make(chan LandingPage)
	go l.populateLandingPage(row, res, ch)
	page = <-ch
	return page, nil
}

func (l *LandingPage) New() error {
	ins, err := database.GetStatement("AddLandingPageStmt")
	if err != nil {
		return err
	}
	ins.Bind(l.Name, l.Start.In(UTC), l.End.In(UTC), l.URL, l.NewWindow, l.MenuPosition)
	_, res, err := ins.Exec()
	if err != nil {
		return err
	}
	id := res.InsertId()
	l.ID = int(id)
	return nil
}

func (l *LandingPage) Save() error {
	upd, err := database.GetStatement("UpdateLandingPageStmt")
	if err != nil {
		return err
	}
	upd.Bind(l.Name, l.Start.In(UTC), l.End.In(UTC), l.URL, l.Content, l.LinkClasses, l.ConversionID, l.ConversionLabel, l.NewWindow, l.MenuPosition, l.ID)
	_, _, err = upd.Exec()
	return err
}

func (l LandingPage) Delete() bool {
	del, err := database.GetStatement("DeleteLandingPageStmt")
	if err != nil {
		return false
	}
	del.Bind(l.ID)
	_, _, err = del.Exec()
	if err != nil {
		return false
	}
	return true
}

func (d LandingPageData) Save() []LandingPageData {
	page := LandingPage{ID: d.LandingPageID}
	data := make([]LandingPageData, 0)
	ch := make(chan []LandingPageData)
	ins, err := database.GetStatement("AddLandingPageDataStmt")
	if err != nil {
		return data
	}
	ins.Bind(d.LandingPageID, d.Key, d.Value)
	ins.Exec()
	go page.getData(ch)
	data = <-ch
	return data
}

func (d LandingPageData) Delete() bool {
	del, err := database.GetStatement("DeleteLandingPageDataStmt")
	if err != nil {
		return false
	}
	del.Bind(d.ID)
	del.Exec()
	return true
}

// Data population

func (l LandingPage) populateLandingPage(row mysql.Row, res mysql.Result, ch chan LandingPage) {
	page := LandingPage{
		ID:              row.Int(res.Map("id")),
		Name:            row.Str(res.Map("name")),
		Start:           row.Time(res.Map("startDate"), UTC),
		End:             row.Time(res.Map("endDate"), UTC),
		URL:             row.Str(res.Map("url")),
		Content:         row.Str(res.Map("pageContent")),
		LinkClasses:     row.Str(res.Map("linkClasses")),
		ConversionID:    row.Str(res.Map("conversionID")),
		ConversionLabel: row.Str(res.Map("conversionLabel")),
		NewWindow:       row.Bool(res.Map("newWindow")),
		MenuPosition:    row.Str(res.Map("menuPosition")),
	}
	dch := make(chan []LandingPageData)
	ich := make(chan LandingPageImages)
	go page.getImages(ich)
	go page.getData(dch)
	page.Images = <-ich
	page.Data = <-dch

	ch <- page
}

func (l LandingPage) getImages(ch chan LandingPageImages) {
	var images LandingPageImages
	sel, err := database.GetStatement("GetLandingPageImagesStmt")
	if err != nil {
		log.Println(err)
		ch <- images
		return
	}
	sel.Bind(l.ID)
	rows, res, err := sel.Exec()
	if err != nil {
		log.Println(err)
		ch <- images
		return
	}
	ich := make(chan LandingPageImage)
	for _, row := range rows {
		go l.populateLandingPageImage(row, res, ich)
	}
	for _, _ = range rows {
		images = append(images, <-ich)
	}
	images.Sort()
	ch <- images
}

func (l LandingPage) populateLandingPageImage(row mysql.Row, res mysql.Result, ch chan LandingPageImage) {
	image := LandingPageImage{
		ID:            row.Int(res.Map("id")),
		LandingPageID: row.Int(res.Map("landingPageID")),
		URL:           row.Str(res.Map("url")),
		Sort:          row.Int(res.Map("sort")),
	}
	ch <- image
}

func (l LandingPage) getData(ch chan []LandingPageData) {
	data := make([]LandingPageData, 0)
	sel, err := database.GetStatement("GetLandingPageDataStmt")
	if err != nil {
		log.Println(err)
		ch <- data
		return
	}
	sel.Bind(l.ID)
	rows, res, err := sel.Exec()
	if err != nil {
		log.Println(err)
		ch <- data
		return
	}
	dch := make(chan LandingPageData)
	for _, row := range rows {
		go l.populateLandingPageData(row, res, dch)
	}
	for _, _ = range rows {
		data = append(data, <-dch)
	}
	ch <- data
}

func (l LandingPage) populateLandingPageData(row mysql.Row, res mysql.Result, ch chan LandingPageData) {
	data := LandingPageData{
		ID:            row.Int(res.Map("id")),
		LandingPageID: row.Int(res.Map("landingPageID")),
		Key:           row.Str(res.Map("dataKey")),
		Value:         row.Str(res.Map("dataValue")),
	}
	ch <- data
}
