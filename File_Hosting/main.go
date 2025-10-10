package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var fileNames = make(map[string]string)

func main() {
	os.MkdirAll("./uploads", 0755)
	http.HandleFunc("/", startHandler)
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/get/", downloadHandler)
	fmt.Println("Start server on port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error start server", err)
	}
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("./templates/upload.html"))
	if r.URL.Query().Get("arc") == "" || r.URL.Query().Get("arc") == "upload" {
		tmpl.Execute(w, nil)
	}
}

func generateShortID(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST запрос", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		panic(err)
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	originalEXT := strings.ToLower(filepath.Ext(handler.Filename))

	fileId, _ := generateShortID(6)
	fmt.Printf("Ваш файл с названием %s был получен, его уникальное id = %s!\n", handler.Filename, fileId)

	fileNames[fileId] = handler.Filename

	dst, err := os.Create("./uploads/" + fileId + originalEXT)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Ошибка сохранения файла на сервер", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "applications/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            fileId,
		"original_name": handler.Filename,
		"url":           "http://localhost:8080/get/" + fileId,
		"extension":     originalEXT,
	})
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileId := strings.TrimPrefix(r.URL.Path, "/get/")

	files, err := os.ReadDir("./uploads")
	if err != nil {
		http.Error(w, "FIle not found", http.StatusNotFound)
		return
	}

	var foundFile string
	originalName, exists := fileNames[fileId]
	if !exists {
		http.NotFound(w, r)
		return
	}

	for _, file := range files {
		filename := file.Name()
		if strings.TrimSuffix(filename, filepath.Ext(filename)) == fileId {
			foundFile = filename
			break
		}
	}

	if foundFile == "" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+originalName+"\"")
	w.Header().Set("Content-type", "application/octet-stream")

	http.ServeFile(w, r, "./uploads/"+foundFile)
	fmt.Printf("Файл с оригинальным названием %s, под короткой ссылкой %s был скачан с сервера!\n", originalName, foundFile)

}
