package main

import (
	"biocadGo/config"
	"biocadGo/db/DBImplementation"
	"biocadGo/db/dbAbstract"
	"biocadGo/src/dir"
	"biocadGo/src/processing"
	"biocadGo/src/webServer/htmlServer"
	"flag"
	"github.com/spf13/viper"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var db dbAbstract.Database
	unProcessed := make(chan string, 5)

	flags := parseFlags()
	config.ConfigInit(flags["config"])

	db = dbInit(viper.GetString("host"), viper.GetString("port"), viper.GetString("DBImplement"))

	wg.Add(1)
	path := viper.GetString("path")
	go dir.CheckUnprocessedFiles(path, &wg, unProcessed)

	wg.Add(1)
	go processing.FilesProcessing(unProcessed, &wg, db)

	wg.Add(1)
	server := htmlServer.Server{}
	go server.Init(viper.GetString("web_server_port"), &wg, db)

	wg.Wait()
}

func parseFlags() map[string]string {

	flag.String("config", "dev", "Имя файла для обработки")
	flag.Parse()
	flags := make(map[string]string)

	flag.Visit(func(f *flag.Flag) {
		flags[f.Name] = f.Value.String()
	})
	return flags
}

func dbInit(host, port, DBImplement string) (dataBase dbAbstract.Database) {
	switch DBImplement {
	case "mongodb":
		dataBase = &DBImplementation.MongoDB{}
	default:
		panic("неверный аргумент типа базы данных")
	}
	dataBase.Init(host, port)
	return dataBase

}
