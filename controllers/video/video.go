package video

import (
	"../../helpers/plate"
	"../../helpers/youtube"
	"../../models"
	"log"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl := plate.NewTemplate(w)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	videos, _ := models.Video{}.GetAll()
	ytvideos, _ := youtube.GetAll(page, 10)
	tmpl.FuncMap["showPrev"] = func() bool {
		return page > 1
	}
	tmpl.FuncMap["showNext"] = func() bool {
		return page < ytvideos.Pages
	}
	tmpl.Bag["PageTitle"] = "Videos"
	tmpl.Bag["videos"] = videos
	tmpl.Bag["ytvideos"] = ytvideos
	tmpl.Bag["next"] = page + 1
	tmpl.Bag["prev"] = page - 1

	tmpl.ParseFile("templates/website/navigation.html", false)
	tmpl.ParseFile("templates/video/index.html", false)

	err = tmpl.Display(w)
	if err != nil {
		log.Println(err)
	}
}

func Sort(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	videos := r.Form["video[]"]
	models.Video{}.UpdateSort(videos)
}
