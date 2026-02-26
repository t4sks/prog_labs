package main

import (
	"auiapp/model"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/save", saveHandler)

	fmt.Println("Server: http://localhost:8080/index")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, _, _ := r.FormFile("file")
		defer file.Close()
		data, _ := io.ReadAll(file)
		works := ParsingFile(data)
		tmpl.Execute(w, works)
		return
	}
	tmpl.Execute(w, nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	var updated []model.ProjectWork
	json.NewDecoder(r.Body).Decode(&updated)

	var b bytes.Buffer
	for _, work := range updated {
		fmt.Fprintf(&b, "\"%s\" \"%s\" %s %s\n", work.Name, work.NameOfWork, work.Date.Format("2006.01.02"), work.Type)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(b.Bytes())
}

func ParsingFile(data []byte) []model.ProjectWork {
	var works []model.ProjectWork
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			works = append(works, ParcingLine(line))
		}
	}
	return works
}

func ParcingLine(line string) model.ProjectWork {

	split := regexp.MustCompile(`"[^"]*"|\d{4}\.\d{2}\.\d{2}|\b[a-f]\b`).FindAllString(line, -1)

	res := model.ProjectWork{}
	if len(split) >= 4 {
		res.Name = strings.Trim(split[0], "\"")
		res.NameOfWork = strings.Trim(split[1], "\"")
		res.Date, _ = time.Parse("2006.01.02", split[2])
		res.Type = split[3]
	}
	return res
}
