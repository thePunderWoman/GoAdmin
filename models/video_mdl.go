package models

import (
	"../helpers/database"
	"../helpers/youtube"
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

func (v Video) Add(ytID string) (Video, error) {
	var video Video
	ytvideo, err := youtube.Get(ytID)
	if err == nil {
		sort := v.GetSort() + 1
		ins, err := database.GetStatement("AddVideoStmt")
		if err != nil {
			return video, err
		}
		ins.Raw.Reset()
		ins.Bind(ytvideo.Video.Details.VideoID, time.Now().In(UTC), sort, ytvideo.Video.Title, ytvideo.Video.Details.Description, ytvideo.Video.Details.VideoID, ytvideo.Video.Details.WatchPage, ytvideo.Video.GetScreenshot())
		_, res, err := ins.Exec()
		if err != nil {
			return video, err
		}
		vID := res.InsertId()
		sel, err := database.GetStatement("GetVideoStmt")
		if err != nil {
			return video, err
		}
		sel.Raw.Reset()
		sel.Bind(vID)
		row, res, err := sel.ExecFirst()
		if err != nil {
			return video, err
		}
		ch := make(chan Video)
		go v.PopulateVideo(row, res, ch)
		video = <-ch
	}
	return video, nil
}

func (v Video) Delete() error {
	del, err := database.GetStatement("DeleteVideoStmt")
	if err != nil {
		return err
	}
	del.Bind(v.ID)
	_, _, err = del.Exec()
	go v.ResetSort()
	return err
}

func (v *Video) ResetSort() {
	videos, err := Video{}.GetAll()
	if err != nil && len(videos) > 0 {
		ids := make([]string, 0)
		for _, video := range videos {
			ids = append(ids, strconv.Itoa(video.ID))
		}
		Video{}.UpdateSort(ids)
	}
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
		upd.Raw.Reset()
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

func (v Video) GetSort() int {
	sel, err := database.GetStatement("GetLastVideoSortStmt")
	if err != nil {
		return 0
	}
	row, res, err := sel.ExecFirst()
	if err != nil {
		return 0
	}
	return row.Int(res.Map("sort"))
}

func (v Videos) Len() int           { return len(v) }
func (v Videos) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Videos) Less(i, j int) bool { return v[i].Sort < v[j].Sort }

func (v *Videos) Sort() {
	sort.Sort(v)
}
