package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Note struct {
	Id          int
	Name        string
	TimeCreate  time.Time
	Desсription string
}

var AllNotes = []Note{}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error reading .env file")
	}
	dbStart()
	handlerAll()

}

func dbStart() {
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("error connect db")
		return
	}
	defer db.Close(context.Background())
	fmt.Println("db started")

	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		fmt.Println("err create db", err)
	}

	fmt.Println("succes created")
}

func handlerAll() {
	http.HandleFunc("/", startHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/save_notify", saveNotifyHandler)
	fmt.Println("Start server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("error while start server")
	}
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home_page.html", "templates/header.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Error connect to db in home")
		return
	}
	defer db.Close(context.Background())

	data, err := db.Query(context.Background(), `SELECT id, name, description, created_at FROM notes`)
	if err != nil {
		fmt.Println("error select all notes ", err)
		return
	}
	AllNotes = []Note{}
	for data.Next() {
		var note Note
		err = data.Scan(&note.Id, &note.Name, &note.Desсription, &note.TimeCreate)
		if err != nil {
			fmt.Println("error select one note ", err)
			return
		}
		AllNotes = append(AllNotes, note)
	}
	tmpl.ExecuteTemplate(w, "home", AllNotes)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	tmpl.ExecuteTemplate(w, "create", nil)
}

func saveNotifyHandler(w http.ResponseWriter, r *http.Request) {
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("error connect in saveFunc")
	} else {
		fmt.Println("succes connect in insert")
	}

	defer db.Close(context.Background())

	title := r.FormValue("name")
	description := r.FormValue("description")

	var insertedId int
	err = db.QueryRow(context.Background(),
		`
		INSERT INTO notes(name, description) VALUES($1, $2) RETURNING id
	`,
		title, description,
	).Scan(&insertedId)
	if err != nil {
		fmt.Println("Error insert notes with id = ", insertedId, err)
	} else {
		fmt.Println("Success insert notes with id = ", insertedId)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
