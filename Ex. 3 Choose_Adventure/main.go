package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

type Chapter struct {
	Title   string
	Story   []string
	Options []Option
}

type Option struct {
	Text string
	Arc  string // КУда перейти при выборе
}

func forestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home page")
}

func introHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Intro page")
}

func loadAllChapter(pathfile string) {
	data, err := os.ReadFile(pathfile)
	if err != nil {
		fmt.Println("error read json", err)
		return
	}

	err = json.Unmarshal(data, &AllchaptersGLOBAL)
	if err != nil {
		fmt.Println("error path url", err)
		return
	}
}

func getChapter(chapterName string) Chapter {
	if chapter, exists := AllchaptersGLOBAL[chapterName]; exists {
		return chapter
	}
	return Chapter{
		Title: "Глава не найдена",
		Story: []string{"Извините, эта глава не существует."},
	}
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/example.html"))
	if r.URL.Query().Get("arc") == "" || r.URL.Query().Get("arc") == "intro" {
		chapter := getChapter("intro")
		tmpl.Execute(w, chapter)
	} else {
		chapter := getChapter(r.URL.Query().Get("arc"))
		tmpl.Execute(w, chapter)
	}

}

var AllchaptersGLOBAL map[string]Chapter

func main() {
	loadAllChapter("story.json")
	http.HandleFunc("/", startHandler)
	fmt.Println("server start on 8080 port")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("error server", err)
	}
}
