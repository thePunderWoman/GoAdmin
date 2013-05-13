package models

import (
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	_ "log"
	"sort"
	"strconv"
	"time"
)

type Posts []Posts
type Post struct {
	ID           int
	Title        string
	Slug         string
	Content      string
	Published    time.Time
	Created      time.Time
	lastModified time.Time
	UserID       int
	MetaTitle    string
	MetaDesc     string
	Keywords     string
	Active       bool
	Categories   Categories
	Comments     Comments
	Author       User
}

type BlogCategories []Category
type BlogCategory struct {
	ID     int
	Name   string
	Slug   string
	Active bool
}

type Comments []Comment
type Comment struct {
	ID       int
	PostID   int
	Name     string
	Email    string
	Comment  string
	Created  time.Time
	Approved bool
	Active   bool
}

func (p Post) GetAll() (Posts, error) {
	var posts Posts
	sel, err := database.GetStatement("GetAllPostsStmt")
	if err != nil {
		return posts, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return posts, err
	}
	ch := make(chan Post)
	for _, row := range rows {
		go p.PopulatePost(row, res, ch)
	}
	for _, _ = range rows {
		posts = append(posts, <-ch)
	}
	return posts, nil
}

func (p Post) PopulatePost(row mysql.Row, res mysql.Result, ch chan Post) {
	post := Post{
		ID:           row.Int(res.Map("blogPostID")),
		Title:        row.Str(res.Map("post_title")),
		Slug:         row.Str(res.Map("slug")),
		Content:      row.Str(res.Map("post_text")),
		Published:    row.Time(res.Map("publishedDate"), UTC),
		Created:      row.Time(res.Map("createdDate"), UTC),
		LastModified: row.Time(res.Map("lastModified"), UTC),
		UserID:       row.Int(res.Map("userID")),
		MetaTitle:    row.Str(res.Map("meta_title")),
		MetaDesc:     row.Str(res.Map("meta_description")),
		Keywords:     row.Str(res.Map("keywords")),
		Active:       row.Bool(res.Map("active")),
	}
	catchan := make(chan BlogCategories)
	comchan := make(chan Comments)
	authchan := make(chan User)
	go func(ch User) {
		author, _ := GetUserByID(post.UserID)
		ch <- author
	}(authchan)
	go func(ch Comments) {
		comments, _ := p.GetComments()
		ch <- comments
	}(comchan)
	go func(ch BlogCategories) {
		categories, _ := p.GetCategories()
		ch <- categories
	}(catchan)

	post.Author = <-authchan
	post.Categories = <-catchan
	post.Comments = <-comchan
	ch <- post
}

func (p Post) GetCategories() (BlogCategories, error) {
	var categories BlogCategories
	sel, err := database.GetStatement("GetPostCategoriesStmt")
	if err != nil {
		return categories, err
	}
	sel.Bind(p.ID)
	rows, res, err := sel.Exec()
	if err != nil {
		return categories, err
	}
	ch := make(chan BlogCategory)
	for _, row := range rows {
		go BlogCategory{}.PopulateCategory(row, res, ch)
	}
	for _, _ = range rows {
		categories = append(categories, <-ch)
	}
	categories.Sort()
	return categories, nil
}

func (c BlogCategory) PopulateCategory(row mysql.Row, res mysql.Result, ch chan BlogCategory) {
	category := BlogCategory{
		ID:     row.Int(res.Map("blogCategoryID")),
		Name:   row.Str(res.Map("name")),
		Slug:   row.Str(res.Map("slug")),
		Active: row.Bool(res.Map("active")),
	}
	ch <- category
}

func (p Post) GetComments() (Comments, error) {
	var comments Comments
	sel, err := database.GetStatement("GetPostCategoriesStmt")
	if err != nil {
		return comments, err
	}
	sel.Bind(p.ID)
	rows, res, err := sel.Exec()
	if err != nil {
		return comments, err
	}
	ch := make(chan Comment)
	for _, row := range rows {
		go Comment{}.PopulateComment(row, res, ch)
	}
	for _, _ = range rows {
		comments = append(comments, <-ch)
	}
	comments.Sort()
	return comments, nil
}

func (c Comment) PopulateComment(row mysql.Row, res mysql.Result, ch chan Comment) {
	comment := Comment{
		ID:       row.Int(res.Map("commentID")),
		PostID:   row.Int(res.Map("blogPostID")),
		Name:     row.Str(res.Map("name")),
		Email:    row.Str(res.Map("email")),
		Comment:  row.Str(res.Map("comment_text")),
		Created:  row.Time(res.Map("createdDate"), UTC),
		Approved: row.Bool(res.Map("approved")),
		Active:   row.Bool(res.Map("active")),
	}
	ch <- comment
}

func (c BlogCategories) Len() int           { return len(c) }
func (c BlogCategories) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c BlogCategories) Less(i, j int) bool { return c[i].Name < (c[j].Name) }

func (c *BlogCategories) Sort() {
	sort.Sort(c)
}

func (c Comments) Len() int           { return len(c) }
func (c Comments) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Comments) Less(i, j int) bool { return c[i].Created.Before(c[j].Created) }

func (c *Comments) Sort() {
	sort.Sort(c)
}

func (p Posts) Len() int           { return len(p) }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Posts) Less(i, j int) bool { return p[i].Published.Before(p[j].Published) }

func (p *Posts) Sort() {
	sort.Sort(p)
}
