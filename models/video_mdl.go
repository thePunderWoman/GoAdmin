package models

import (
	"../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	"log"
	"sort"
	"strconv"
	"time"
)

type Video struct {
	ID          int
	EmbedLink   string
	DateAdded   time.Time
	Sort        int
	Title       string
	Description string
	YouTubeID   string
	WatchPage   string
	Screenshot  string
}

type Videos []Video

func (v Video) GetAll() (Videos, error) {
	var videos Videos
	sel, err := database.GetStatement("GetAllVideosStmt")
	if err != nil {
		return videos, err
	}
	rows, res, err := sel.Exec()
	if err != nil {
		return videos, err
	}
	ch := make(chan Video)
	for _, row := range rows {
		go v.PopulateVideo(row, res, ch)
	}
	for _, _ = range rows {
		videos = append(videos, <-ch)
	}
	videos.Sort()
	return videos, nil
}

func (v Video) UpdateSort(videos []string) {
	sort := 0
	for _, video := range videos {
		sort += 1
		upd, err := database.GetStatement("UpdateVideoSortStmt")
		if err != nil {
			log.Println(err)
			return
		}
		videoID, _ := strconv.Atoi(video)
		upd.Reset()
		upd.Bind(sort, videoID)
		upd.Exec()
	}
}

func (v Video) PopulateVideo(row mysql.Row, res mysql.Result, ch chan Video) {
	video := Video{
		ID:          row.Int(res.Map("videoID")),
		EmbedLink:   row.Str(res.Map("embed_link")),
		DateAdded:   row.Time(res.Map("dateAdded"), UTC),
		Sort:        row.Int(res.Map("sort")),
		Title:       row.Str(res.Map("title")),
		Description: row.Str(res.Map("description")),
		YouTubeID:   row.Str(res.Map("youtubeID")),
		WatchPage:   row.Str(res.Map("watchpage")),
		Screenshot:  row.Str(res.Map("screenshot")),
	}
	ch <- video
}

func (v Videos) Len() int           { return len(v) }
func (v Videos) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Videos) Less(i, j int) bool { return v[i].Sort < v[j].Sort }

func (v *Videos) Sort() {
	sort.Sort(v)
}
