package webServer

import (
	"Biocad2/db/dbInterface"
	"Biocad2/src/message"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"sync"
)

type Server struct {
	ServerDB   dbInterface.DBI
	mainConfig *viper.Viper
}

func (s *Server) Init(wg *sync.WaitGroup, db dbInterface.DBI, conf *viper.Viper) {
	s.ServerDB = db
	s.mainConfig = conf
	defer wg.Done()
	mux := http.NewServeMux()
	mux.HandleFunc("/index", s.IndexHandler)
	mux.HandleFunc("/json", s.JsonHandler)
	err := http.ListenAndServe(s.mainConfig.GetString("web_server_port"), mux)
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
