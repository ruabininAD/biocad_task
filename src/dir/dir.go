package dir

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"

	"sync"
	"time"
)

func CheckUnprocessedFiles(mainConfig *viper.Viper, wg *sync.WaitGroup, unprocessed chan string) {
	defer wg.Done()
	directory := mainConfig.GetString("path")
	for {
		checkUnprocessedFiles(directory, unprocessed)

		// Периодичность проверок
		timeout := mainConfig.GetInt("frequency_of_dir_checks")
		time.Sleep(time.Second * time.Duration(timeout))
	}

}

// Проверка файла на не папку, отсутствие префикса Processed, наличие tsv расширение
// Отправка файла в канал unprocessed
func checkUnprocessedFiles(directory string, unprocessed chan string) {

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Println("Ошибка при чтении директории:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && !isFileInProcessing(file.Name()) && !isFileProcessed(file.Name()) {
			newPath, _ := inProcessFileName(directory + "/" + file.Name())
			log.Println(file.Name() + " отмечена как in process")
			unprocessed <- newPath
		}
	}
}

func inProcessFileName(filePath string) (newPath string, err error) {
	directory, fileName := filepath.Split(filePath)
	newFilename := fmt.Sprintf("in_procss_%s", fileName)
	newPath = filepath.Join(directory, newFilename)

	err = os.Rename(filePath, newPath)
	if err != nil {
		return newPath, err
	}

	return newPath, nil
}

func isFileInProcessing(filename string) bool {
	return len(filename) >= 10 && filename[:10] == "in_procss_"
}
func isFileProcessed(filename string) bool {
	return len(filename) >= 10 && filename[:10] == "processed_"
}
