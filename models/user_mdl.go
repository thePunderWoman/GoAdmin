package models

import (
	"../helpers/database"
	"bytes"
	"crypto/md5"
	"errors"
	_ "log"
	"time"
)

var (
	UTC, _ = time.LoadLocation("UTC")
)

type User struct {
	ID        int
	Username  string
	Email     string
	Fname     string
	Lname     string
	DateAdded time.Time
	IsActive  bool
	SuperUser bool
	Biography string
	Photo     string
	Modules   []Module
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

	epw, err := Md5Encrypt(password)
	if err != nil {
		return user, err
	}

	sel.Bind(username, epw)

	row, res, err := sel.ExecFirst()
	if database.MysqlError(err) {
		return user, err
	}

	id := res.Map("id")
	uname := res.Map("username")
	fname := res.Map("fname")
	lname := res.Map("lname")
	super := res.Map("superUser")

	if err != nil { // Must be something wrong with the db, lets bail
		return user, err
	} else if row != nil { // populate history object
		user = User{
			ID:        row.Int(id),
			Username:  row.Str(uname),
			Fname:     row.Str(fname),
			Lname:     row.Str(lname),
			SuperUser: row.Bool(super),
		}
	}

	return user, err
}

func GetUser(id int) (u User, err error) {
	sel, err := database.GetStatement("getUserByIDStmt")
	if err != nil {
		return u, err
	}

	sel.Bind(id)

	row, res, err := sel.ExecFirst()
	if database.MysqlError(err) {
		return u, err
	}

	idval := res.Map("id")
	uname := res.Map("username")
	email := res.Map("email")
	fname := res.Map("fname")
	lname := res.Map("lname")
	dateAdded := res.Map("dateAdded")
	active := res.Map("isActive")
	super := res.Map("superUser")
	bio := res.Map("biography")
	photo := res.Map("photo")

	if err != nil { // Must be something wrong with the db, lets bail
		return u, err
	} else if row != nil { // populate history object
		u = User{
			ID:        row.Int(idval),
			Username:  row.Str(uname),
			Email:     row.Str(email),
			Fname:     row.Str(fname),
			Lname:     row.Str(lname),
			DateAdded: row.Time(dateAdded, UTC),
			IsActive:  row.Bool(active),
			SuperUser: row.Bool(super),
			Biography: row.Str(bio),
			Photo:     row.Str(photo),
		}
		err = u.GetModules()
		if err != nil {
			return u, err
		}
	}

	return u, err

}

func (u *User) GetModules() error {
	sel, err := database.GetStatement("userModulesStmt")
	if err != nil {
		return err
	}

	sel.Bind(u.ID)

	rows, res, err := sel.Exec()
	if database.MysqlError(err) {
		return err
	}

	id := res.Map("id")
	mod := res.Map("module")
	modPath := res.Map("module_path")
	imgPath := res.Map("img_path")

	var modules []Module

	for _, row := range rows {
		m := Module{
			ID:          row.Int(id),
			Module:      row.Str(mod),
			Module_path: row.Str(modPath),
			Img_path:    row.Str(imgPath),
		}
		modules = append(modules, m)
	}
	u.Modules = modules
	return nil
}

func Md5Encrypt(str string) (string, error) {
	if str == "" {
		return "", errors.New("Invalid string parameter")
	}
	h := md5.New()
	h.Write([]byte(str))
	var buf bytes.Buffer
	_, err := buf.Write(h.Sum(nil))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
