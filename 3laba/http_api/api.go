package http_api

import (
	"auiapp/function"
	"auiapp/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

var tmpl = template.Must(template.ParseFiles("web/templates/index.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, _, _ := r.FormFile("file")
		defer file.Close()
		data, _ := io.ReadAll(file)
		works := function.ParsingFile(data)
		tmpl.Execute(w, works)
		return
	}
	tmpl.Execute(w, nil)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	var updated []model.ProjectWork
	json.NewDecoder(r.Body).Decode(&updated)
	var b bytes.Buffer
	for _, work := range updated {
		line := fmt.Sprintf(
			"\"%s\" \"%s\" %s %s",
			work.Name,
			work.NameOfWork,
			work.Date.Format("2006.01.02"),
			work.Type,
		)
		parsed := function.ParcingLine(line)
		if parsed.Name == "" {
			continue
		}
		fmt.Fprintln(&b, line)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(b.Bytes())
}
