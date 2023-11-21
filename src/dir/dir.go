package dir

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
	"time"
)

func CheckUnprocessedFiles(directory string, wg *sync.WaitGroup, unprocessed chan string) {
	defer wg.Done()
	for {
		checkUnprocessedFiles(directory, unprocessed)

		// Периодичность проверок
		timeout := viper.GetInt("frequency_of_dir_checks")
		time.Sleep(time.Second * time.Duration(timeout))
	}

}

// Проверка файла на не папку, отсутствие префикса Processed, наличие tsv расширение
// Отправка файла в канал unprocessed
func checkUnprocessedFiles(directory string, unprocessed chan string) {

	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Ошибка при чтении директории:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && !isFileProcessed(file.Name()) && isTSVFile(file.Name()) {
			unprocessed <- viper.GetString("path") + "/" + file.Name()
		}
	}
}

func isFileProcessed(filename string) bool {
	return len(filename) >= 10 && filename[:10] == "processed_"
}

func isTSVFile(filename string) bool {
	return strings.HasSuffix(filename, ".tsv")
}
