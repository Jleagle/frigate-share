package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var caser = cases.Title(language.BritishEnglish)

func main() {

	var templates = template.Must(template.ParseFiles("main.gohtml"))

	http.HandleFunc("GET /logo.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "logo.jpg")
	})

	http.HandleFunc("GET /manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "manifest.json")
	})

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

		if len(events) > 50 {
			events = events[:50]
		}

		err = templates.ExecuteTemplate(w, "main.gohtml", Data{Events: events})
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	http.HandleFunc("GET /events/{id}", func(w http.ResponseWriter, r *http.Request) {

		url := fmt.Sprintf("http://frigate:5000/api/events/%s/clip.mp4", r.PathValue("id"))

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating origin request: %v", err), http.StatusInternalServerError)
			return
		}

		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" {
			req.Header.Set("Range", rangeHeader)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching video from origin: %v", err), http.StatusBadGateway)
			return
		}

		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)

		_, err = io.Copy(w, resp.Body)
		if err != nil && !strings.Contains(err.Error(), "broken pipe") {
			fmt.Printf("Error copying response body: %v\n", err)
		}
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
	Camera    string  `json:"camera"`
	ID        string  `json:"id"`
	Label     string  `json:"label"`
	StartTime float64 `json:"start_time"`
}

func (e Event) Start() string {
	return time.Unix(int64(e.StartTime), 0).Format(time.Stamp)
}

func (e Event) LabelNice() string {
	return caser.String(e.Label)
}

func (e Event) CameraNice() string {
	return caser.String(e.Camera)
}

//func (e Event) ScoreNice() int {
//	return int(e.Data.Score * 100)
//}
