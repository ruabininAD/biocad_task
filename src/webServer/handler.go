package webServer

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {

	//  .html в виде байтов
	htmlContent, err := os.ReadFile("config/configFile/index.html") //fixme
	if err != nil {
		http.Error(w, "Unable to load HTML template", http.StatusInternalServerError)
		return
	}

	t, err := template.New("index").Funcs(template.FuncMap{"add": add}).Parse(string(htmlContent))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		data := PageData{Title: "Search Page"}
		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {

		id := r.FormValue("query")
		pageSize := s.mainConfig.GetInt("pageSize")
		pageNumber, _ := strconv.Atoi(r.FormValue("page"))

		mess, total := s.ServerDB.GetById(id, pageNumber, pageSize)

		data := PageData{Title: id, Mes: mess, Total: int(total), PageN: pageNumber, UnitGuid: id, PageSize: pageSize}
		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("запрос id: %s, pageSize:%v, pageNumber%v \n", id, pageSize, pageNumber)
	}

}

// Ваша функция обработки запросов
func (s *Server) JsonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		queryValues := r.URL.Query()
		id := queryValues.Get("id")
		pageSize := s.mainConfig.GetInt("pageSize")
		pageNumber, _ := strconv.Atoi(queryValues.Get("page"))

		mess, total := s.ServerDB.GetById(id, pageNumber, pageSize)

		data := PageData{Title: id, Mes: mess, Total: int(total), PageN: pageNumber, UnitGuid: id, PageSize: pageSize}

		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// func sum page for html page
func add(a, b int, Total, PageSize int) int {
	res := a + b
	if res < 1 {
		return 1
	}

	result := Total % PageSize
	if result != 0 {
		result = 1
	}
	pageCount := Total/PageSize + result
	if res > pageCount {
		return 1
	}

	return a + b
}
