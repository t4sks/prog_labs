package http_api

import (
	"auiapp/function"
	"auiapp/model"
	"encoding/json"
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
	b := function.ObjectTobytes(updated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(b)
}

func ApplyCommandsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Ошибка чтения формы", http.StatusInternalServerError)
		return
	}

	data := r.FormValue("data")
	if data == "" {
		http.Error(w, "Ошибка данных", http.StatusInternalServerError)
		return
	}
	var works []model.ProjectWork
	err = json.Unmarshal([]byte(data), &works)
	if err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusInternalServerError)
		return
	}
	file, _, err := r.FormFile("commands")
	if err != nil {
		http.Error(w, "Файл с командами не передан", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	commandBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка чтения команд", http.StatusInternalServerError)
		return
	}
	works = function.ReadExecCommandFile(works, commandBytes)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(works)
}
