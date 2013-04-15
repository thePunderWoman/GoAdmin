package models

import (
	"../helpers/database"
	"crypto/md5"
	"log"
	"time"
)

var (
	UTC, _ = time.LoadLocation("UTC")
)

type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	Fname     string
	Lname     string
	DateAdded time.Time
	IsActive  bool
	SuperUser bool
	Biography string
	Photo     string
}

type Module struct {
	ID          int
	Module      string
	Module_path string
	Img_path    string
}

func Authenticate(username string, password string) (user User, err error) {
	sel, err := database.GetStatement("authenticateUserStmt")
	if err != nil {
		return user, err
	}

	sel.Bind(username, password)

	row, res, err := sel.ExecFirst()
	if database.MysqlError(err) {
		return user, err
	}

	id := res.Map("id")
	uname := res.Map("username")
	pword := res.Map("password")
	email := res.Map("email")
	fname := res.Map("fname")
	lname := res.Map("lname")
	dateAdded := res.Map("dateAdded")
	active := res.Map("isActive")
	super := res.Map("superUser")
	bio := res.Map("biography")
	photo := res.Map("photo")

	if err != nil { // Must be something wrong with the db, lets bail
		return user, err
	} else if row != nil { // populate history object
		user = User{
			ID:        row.Int(id),
			Username:  row.Str(uname),
			Password:  row.Str(pword),
			Email:     row.Str(email),
			Fname:     row.Str(fname),
			Lname:     row.Str(lname),
			DateAdded: row.Time(dateAdded, UTC),
			IsActive:  row.Bool(active),
			SuperUser: row.Bool(super),
			Biography: row.Str(bio),
			Photo:     row.Str(photo),
		}
	}

	return user, err
}

func EncryptAll() {
	h := md5.New()
	sel, err := database.GetStatement("allUserStmt")
	if err != nil {
		log.Println(err)
		return
	}

	rows, res, err := sel.Exec()
	if database.MysqlError(err) {
		log.Println(err)
		return
	}

	//id := res.Map("id")
	uname := res.Map("username")
	pword := res.Map("password")

	if err != nil {
		log.Println(err)
		return
	} else if len(rows) > 0 {
		for _, row := range rows {
			username := row.Str(uname)
			log.Println(username)
			pw := row.Str(pword)
			log.Println(pw)
			h.Write([]byte(pw))
			epw := h.Sum(nil)
			log.Println(epw)
		}
	}
}
