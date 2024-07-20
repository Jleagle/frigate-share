package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var caser = cases.Title(language.BritishEnglish)

func main() {

	var templates = template.Must(template.ParseFiles("main.gohtml"))

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		resp, err := http.Get("http://frigate:5000/api/events")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		//b, err := os.ReadFile("events.json")
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}

		var events []Event
		err = json.Unmarshal(b, &events)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = templates.ExecuteTemplate(w, "main.gohtml", Data{Events: events})
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	http.HandleFunc("GET /events/{id}", func(w http.ResponseWriter, r *http.Request) {

		resp, err := http.Get(fmt.Sprintf("http://frigate:5000/api/events/%s/clip.mp4", r.PathValue("id")))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(b))
	})

	err := http.ListenAndServe("0.0.0.0:5002", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Data struct {
	Events []Event
}

type Event struct {
	Box    interface{} `json:"box"`
	Camera string      `json:"camera"`
	Data   struct {
		Attributes []interface{} `json:"attributes"`
		Box        []float64     `json:"box"`
		Region     []float64     `json:"region"`
		Score      float64       `json:"score"`
		TopScore   float64       `json:"top_score"`
		Type       string        `json:"type"`
	} `json:"data"`
	EndTime            float64       `json:"end_time"`
	FalsePositive      interface{}   `json:"false_positive"`
	HasClip            bool          `json:"has_clip"`
	HasSnapshot        bool          `json:"has_snapshot"`
	ID                 string        `json:"id"`
	Label              string        `json:"label"`
	PlusID             interface{}   `json:"plus_id"`
	RetainIndefinitely bool          `json:"retain_indefinitely"`
	StartTime          float64       `json:"start_time"`
	SubLabel           interface{}   `json:"sub_label"`
	Thumbnail          string        `json:"thumbnail"`
	TopScore           interface{}   `json:"top_score"`
	Zones              []interface{} `json:"zones"`
}

func (e Event) Start() string { //
	return time.Unix(int64(e.StartTime), 0).Format(time.Stamp)
}

func (e Event) LabelNice() string {
	return caser.String(e.Label)
}

func (e Event) CameraNice() string {
	return caser.String(e.Camera)
}

func (e Event) ScoreNice() int {
	return int(e.Data.Score * 100)
}

func (e Event) MP4() string {
	return fmt.Sprintf("/events/%s.mp4", e.ID)
}
