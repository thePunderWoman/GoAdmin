package models

import (
	"../helpers/database"
	"../helpers/email"
	"bytes"
	"crypto/md5"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"strings"
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

func (u *User) New() error {
	// new user
	// check if username exists
	uchan := make(chan bool)
	echan := make(chan bool)
	go func(username string) {
		_, err := GetUserByUsername(username)
		if err != nil {
			uchan <- false
		} else {
			uchan <- true
		}
	}(u.Username)

	go func(email string) {
		// check if email exists
		_, err := GetUserByEmail(email)
		if err != nil {
			echan <- false
		} else {
			echan <- true
		}
	}(u.Email)

	uexists := <-uchan
	eexists := <-echan
	if uexists {
		return errors.New("A User account with that username already exists.")
	}
	if eexists {
		return errors.New("A User account with that email already exists.")
	}
	// add user
	ins, err := database.GetStatement("registerUserStmt")
	if err != nil {
		return err
	}

	params := struct {
		Username string
		Email    string
		Fname    string
		Lname    string
	}{}

	params.Username = u.Username
	params.Email = u.Email
	params.Fname = u.Fname
	params.Lname = u.Lname

	ins.Bind(&params)

	_, _, err = ins.Exec()
	if err != nil {
		return err
	}

	sel, err := database.GetStatement("getID")
	if err != nil {
		return err
	}

	row, res, err := sel.ExecFirst()
	if err != nil {
		return err
	}

	id := res.Map("id")
	u.ID = row.Int(id)
	err = u.ResetPassword()
	if err != nil {
		return err
	}
	u.SendNewUserEmail()

	return nil
}

func (u *User) Save() error {
	// new user
	// check if username exists
	if u.ID == 0 {
		uchan := make(chan bool)
		echan := make(chan bool)
		go func(username string) {
			_, err := GetUserByUsername(username)
			if err != nil {
				uchan <- false
			} else {
				uchan <- true
			}
		}(u.Username)

		go func(email string) {
			// check if email exists
			_, err := GetUserByEmail(email)
			if err != nil {
				echan <- false
			} else {
				echan <- true
			}
		}(u.Email)

		uexists := <-uchan
		eexists := <-echan
		if uexists {
			return errors.New("A User account with that username already exists.")
		}
		if eexists {
			return errors.New("A User account with that email already exists.")
		}
		// add user
		ins, err := database.GetStatement("addUserStmt")
		if err != nil {
			return err
		}

		params := struct {
			Username  string
			Email     string
			Fname     string
			Lname     string
			Biography string
			Photo     string
			IsActive  bool
			SuperUser bool
		}{}

		params.Username = u.Username
		params.Email = u.Email
		params.Fname = u.Fname
		params.Lname = u.Lname
		params.Biography = u.Biography
		params.Photo = u.Photo
		params.IsActive = u.IsActive
		params.SuperUser = u.SuperUser

		ins.Bind(&params)

		_, _, err = ins.Exec()
		if err != nil {
			return err
		}

		sel, err := database.GetStatement("getID")
		if err != nil {
			return err
		}

		row, res, err := sel.ExecFirst()
		if err != nil {
			return err
		}

		id := res.Map("id")
		u.ID = row.Int(id)

		u.SavePassword()
	} else {
		// add user
		upd, err := database.GetStatement("updateUserStmt")
		if err != nil {
			return err
		}

		params := struct {
			Username  string
			Email     string
			Fname     string
			Lname     string
			Biography string
			Photo     string
			IsActive  bool
			SuperUser bool
			ID        int
		}{}

		params.Username = u.Username
		params.Email = u.Email
		params.Fname = u.Fname
		params.Lname = u.Lname
		params.Biography = u.Biography
		params.Photo = u.Photo
		params.IsActive = u.IsActive
		params.SuperUser = u.SuperUser
		params.ID = u.ID

		upd.Bind(&params)

		_, _, err = upd.Exec()
		if err != nil {
			return err
		}
		if strings.TrimSpace(u.Password) != "" {
			u.SavePassword()
		}
	}
	return nil
}

func GetUserByID(id int) (u User, err error) {
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
	if u.ID == 0 {
		return u, errors.New("User not found")
	}

	return u, err
}

func GetUserByUsername(username string) (u User, err error) {
	sel, err := database.GetStatement("getUserByUsernameStmt")
	if err != nil {
		return u, err
	}

	sel.Bind(username)

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
	if u.ID == 0 {
		return u, errors.New("User not found")
	}

	return u, err
}

func GetUserByEmail(email string) (u User, err error) {
	sel, err := database.GetStatement("getUserByEmailStmt")
	if err != nil {
		return u, err
	}

	sel.Bind(email)

	row, res, err := sel.ExecFirst()
	if database.MysqlError(err) {
		return u, err
	}

	idval := res.Map("id")
	uname := res.Map("username")
	emailaddr := res.Map("email")
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
			Email:     row.Str(emailaddr),
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
	if u.ID == 0 {
		return u, errors.New("User not found")
	}

	return u, err
}

func (u *User) GetAll() (users []User, err error) {
	sel, err := database.GetStatement("getAllUserStmt")
	if err != nil {
		return users, err
	}
	rows, res, err := sel.Exec()
	if database.MysqlError(err) {
		return users, err
	}

	idval := res.Map("id")
	uname := res.Map("username")
	emailaddr := res.Map("email")
	fname := res.Map("fname")
	lname := res.Map("lname")
	dateAdded := res.Map("dateAdded")
	active := res.Map("isActive")
	super := res.Map("superUser")
	bio := res.Map("biography")
	photo := res.Map("photo")

	for _, row := range rows {
		u := User{
			ID:        row.Int(idval),
			Username:  row.Str(uname),
			Email:     row.Str(emailaddr),
			Fname:     row.Str(fname),
			Lname:     row.Str(lname),
			DateAdded: row.Time(dateAdded, UTC),
			IsActive:  row.Bool(active),
			SuperUser: row.Bool(super),
			Biography: row.Str(bio),
			Photo:     row.Str(photo),
		}
		users = append(users, u)
	}

	return users, nil
}

func (u *User) SetStatus(status bool) error {
	stat := 0
	if status {
		stat = 1
	}

	upd, err := database.GetStatement("setUserStatusStmt")
	if err != nil {
		return err
	}

	upd.Bind(stat, u.ID)

	_, _, err = upd.Exec()
	if database.MysqlError(err) {
		return err
	}
	return nil
}

func (u *User) SaveModules(modules []string) {
	// clear user modules
	del, err := database.GetStatement("clearUserModuleStmt")
	if err != nil {
		log.Println(err)
		return
	}

	del.Bind(u.ID)

	_, _, err = del.Exec()
	if database.MysqlError(err) {
		log.Println(err)
		return
	}

	c := make(chan int)

	// re-add user modules
	for _, module := range modules {
		modID, _ := strconv.Atoi(module)
		go u.AddUserModule(modID, c)
	}

	for _, _ = range modules {
		<-c
	}

}

func (u *User) AddUserModule(mID int, ch chan int) {
	ins, err := database.GetStatement("addModuleToUserStmt")
	if err != nil {
		log.Println(err)
		ch <- mID
		return
	}

	ins.Reset()
	ins.Bind(u.ID, mID)

	_, _, err = ins.Exec()
	if database.MysqlError(err) {
		log.Println(err)
		ch <- mID
		return
	}

	ch <- mID
}

func (u *User) Delete() error {
	del, err := database.GetStatement("clearUserModuleStmt")
	if err != nil {
		return err
	}

	del.Bind(u.ID)

	_, _, err = del.Exec()
	if database.MysqlError(err) {
		return err
	}

	del2, err := database.GetStatement("deleteUserStmt")
	if err != nil {
		return err
	}

	del2.Bind(u.ID)

	_, _, err = del2.Exec()
	if database.MysqlError(err) {
		return err
	}
	return nil
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

func GetAllModules() (modules []Module, err error) {
	sel, err := database.GetStatement("getAllModulesStmt")
	if err != nil {
		return modules, err
	}

	rows, res, err := sel.Exec()
	if database.MysqlError(err) {
		return modules, err
	}

	id := res.Map("id")
	mod := res.Map("module")
	modPath := res.Map("module_path")
	imgPath := res.Map("img_path")

	for _, row := range rows {
		m := Module{
			ID:          row.Int(id),
			Module:      row.Str(mod),
			Module_path: row.Str(modPath),
			Img_path:    row.Str(imgPath),
		}
		modules = append(modules, m)
	}

	return modules, nil
}

func (u *User) ResetPassword() error {
	newpassword := GeneratePassword()
	u.SendPasswordEmail(newpassword)

	encpassword, err := Md5Encrypt(newpassword)
	if err != nil {
		return err
	}

	upd, err := database.GetStatement("setUserPasswordStmt")
	if err != nil {
		return err
	}

	upd.Bind(encpassword, u.ID)
	_, _, err = upd.Exec()
	return nil
}

func (u *User) SavePassword() {
	upd, err := database.GetStatement("setUserPasswordStmt")
	if err != nil {
		log.Println(err)
		return
	}

	upd.Bind(u.Password, u.ID)
	_, _, err = upd.Exec()
	return
}

func (u *User) SendPasswordEmail(password string) {
	tos := []string{u.Email}
	body := "<p>Hi " + u.Fname + ",</p>"
	body += "<p>Here is your new password for the <a href=\"http://admin.curtmfg.com\">CURT Administration</a> site.<br /></br />Password: " + password + "</p>"
	body += "<p>If you did not request this password reset, please contact <a href=\"mailto:websupport@curtmfg.com\">Web Support</a>.</p>"
	body += "<p>Thanks,<br />The Ecommerce Developer Team</p>"
	subject := "CURT Administration Password Reset"
	email.Send(tos, subject, body, true)
}

func (u *User) SendNewUserEmail() {
	tos := []string{"websupport@curtmfg.com"}
	body := "<div style='margin-top: 15px;font-family: Arial;font-size: 10pt;'>"
	body += "<div style='border-bottom: 2px solid #999'>"
	body += "<p>A new account has been created with the e-mail {" + u.Email + "}. </p>"
	body += "<p style='margin:2px 0px'>Name: <strong>" + u.Fname + " " + u.Lname + "</strong></p>"
	body += "<p style='margin:2px 0px'>Please login to the admin section of the CURT Administration and activate the account.</p>"
	body += "</div>"
	body += "<br /><span style='color:#999'>Thank you,</span>"
	body += "<br /><br /><br />"
	body += "<span style='line-height:75px;color:#999'>CURT Administration</span>"
	body += "</div>"
	subject := "CURT Administration Account Sign Up"
	email.Send(tos, subject, body, true)
}

func GeneratePassword() string {
	charlist := "ABCDEFGHJKMNOPQRSTUVWXYZabcdefghjkmnopqrstuvwxyz23456789!@#$^&*?"
	charslice := strings.Split(charlist, "")
	targetlength := 8
	newpw := ""
	var index int
	for i := 0; i < targetlength; i++ {
		rand.Seed(time.Now().UnixNano())
		index = rand.Intn(len(charslice))
		newpw += charslice[index]
	}
	return newpw
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
