package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type pathURL struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func readJson(pathfile string) map[string]string {
	data, err := os.ReadFile(pathfile)
	if err != nil {
		fmt.Println("error read json", err)
		return nil
	}

	var pathURLs []pathURL
	err = json.Unmarshal(data, &pathURLs)
	if err != nil {
		fmt.Println("error path url", err)
		return nil
	}

	paths := make(map[string]string)
	for _, pu := range pathURLs {
		paths[pu.Path] = pu.URL
	}
	return paths

}

func writeNewDataJson(pathfile string, data pathURL) {
	dataFromFile, err := os.ReadFile("test.json")
	if err != nil {
		fmt.Printf("not json file")
	}
	urls := []pathURL{}
	json.Unmarshal(dataFromFile, &urls)
	urls = append(urls, data)
	dataToWrite, _ := json.MarshalIndent(urls, "", "	")
	err = os.WriteFile(pathfile, dataToWrite, 0644)
	if err != nil {
		fmt.Println("err write file", err)
	}
}

func createRedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paths := readJson("test.json")
		if target, ok := paths[r.URL.Path]; ok {
			http.Redirect(w, r, target, http.StatusFound)
			return
		}
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "home page")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w,
			`
            <form method="POST">
                <input type="text" name="path" placeholder="Enter your short path">
				<input type="text" name="url" placeholder="Enter your long url">
                <button type="submit">Submit</button>
            </form>
        `)
	} else if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "err parse form %v\n", err)
			return
		}
		var data pathURL
		data.Path = r.FormValue("path")
		if !strings.HasPrefix(data.Path, "/") {
			data.Path = "/" + data.Path
		}
		data.URL = r.FormValue("url")
		writeNewDataJson("test.json", data)
	}
}

func main() {
	for {
		http.HandleFunc("/home", homeHandler)
		http.HandleFunc("/form", formHandler)
		http.HandleFunc("/", createRedirectHandler())
		fmt.Println("start server at 8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("error = ", err)
		}
	}

}
