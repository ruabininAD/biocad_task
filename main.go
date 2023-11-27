package main

import (
	"Biocad2/config"
	"Biocad2/db/dbImplement"
	"Biocad2/db/dbInterface"
	"Biocad2/src/dir"
	"Biocad2/src/processing"
	"Biocad2/src/webServer"
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

func main() {
	fmt.Println("start")

	var wg sync.WaitGroup
	var db dbInterface.DBI
	unProcessed := make(chan string, 5)
	updatePDF := make(chan string, 5)

	flags := config.ParseFlags()
	mainConfig, dbConfig, PDFConfig := config.Init(flags["config"])

	db = dbInit(dbConfig)

	wg.Add(1)
	go dir.CheckUnprocessedFiles(mainConfig, &wg, unProcessed)

	wg.Add(1)
	go processing.FilesProcessing(PDFConfig, unProcessed, updatePDF, &wg, db)

	server := webServer.Server{}
	go server.Init(&wg, db, mainConfig)

	wg.Wait()
}

func dbInit(dbConfig *viper.Viper) (dataBase dbInterface.DBI) {

	switch dbConfig.GetString("DBImplement") {
	case "mongodb":
		dataBase = &dbImplement.MongoDB{}
	default:
		panic("неверный аргумент типа базы данных")
	}
	dataBase.Init(dbConfig)
	return dataBase

}
