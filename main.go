package main

import (
	"biocadGo/config"
	"biocadGo/db"
	"biocadGo/db/DBImplementation"
	"biocadGo/src/dir"
	"biocadGo/src/processing"
	"flag"
	"github.com/spf13/viper"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var db db.Database
	unProcessed := make(chan string, 5)

	flags := parseFlags()
	config.ConfigInit(flags["config"])

	db = dbInit(viper.GetString("host"), viper.GetString("port"), viper.GetString("DBImplement"))

	wg.Add(1)
	path := viper.GetString("path")
	go dir.CheckUnprocessedFiles(path, &wg, unProcessed)

	wg.Add(1)
	go processing.FilesProcessing(unProcessed, &wg, db)

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

func dbInit(host, port, DBImplement string) (dataBase db.Database) {
	switch DBImplement {
	case "mongodb":
		dataBase = &DBImplementation.MongoDB{}
	default:
		panic("не верный аргумент типа базы данных")
	}
	dataBase.Init(host, port)
	return dataBase

}
