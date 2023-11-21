package htmlServer

import (
	"biocadGo/db/dbAbstract"
	"biocadGo/src/message"
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Server struct {
	ServerDB dbAbstract.Database
}

func (s *Server) Init(port string, wg *sync.WaitGroup, db dbAbstract.Database) {
	s.ServerDB = db
	defer wg.Done()
	mux := http.NewServeMux()
	mux.HandleFunc("/index", s.indexHandler)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println(err) // fixme
	}

}

type PageData struct {
	Title    string
	Mes      []message.Message
	Total    int
	PageN    int
	UnitGuid string
	PageSize int //Число сток на странице
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {

	//  .html в виде байтов
	htmlContent, err := os.ReadFile("config/configuration/public/index.html") //fixme
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
		pageSize := viper.GetInt("pageSize")
		pageNumber, _ := strconv.Atoi(r.FormValue("page")) //fixme

		mess, total := s.ServerDB.GetById(id, pageNumber, pageSize)

		data := PageData{Title: id, Mes: mess, Total: int(total), PageN: pageNumber, UnitGuid: id, PageSize: pageSize}
		err = t.Execute(w, data)
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
