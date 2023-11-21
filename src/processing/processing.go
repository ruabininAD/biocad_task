package processing

import (
	"biocadGo/db"
	"biocadGo/src/message"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Retrieves the file paths from "unprocessing" channel
// for each file: parses the file, the data is sent to the database, marks the file as processed,
func FilesProcessing(unProcessing chan string, wg *sync.WaitGroup, db db.Database) {
	defer wg.Done()

	for filePath := range unProcessing {

		newFilePath, _ := processFileName(filePath)

		messages, err := parseFile(newFilePath)
		if err != nil {
			unProcessFileName(filePath)
		}

		err = db.AddFile(messages)
		if err != nil {
			unProcessFileName(filePath)
		}

	}
}

func parseFile(filePath string) ([]message.Message, error) {
	var level, bit, invertBit int

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
	}
	defer file.Close()

	// Создаем ридер для файла
	reader := csv.NewReader(file)
	reader.Comma = '\t' // Устанавливаем символ-разделитель как табуляцию

	// Читаем содержимое файла
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
	}

	// Создаем слайс структур для хранения данных
	var messages []message.Message

	// Проходимся по строкам и заполняем структуру данных
	for numberLine, line := range lines {
		if numberLine == 0 {
			continue
		}

		// Пропускаем заголовок, если он есть
		if strings.HasPrefix(line[0], "ID") {
			continue
		}

		level, bit, invertBit, err = transformDataType(line[8], line[13], line[14])
		if err != nil {
			fmt.Println("Ошибка при приведении типов n|level|bit файл:"+filePath+": ", err)
		}

		currentMessage := message.Message{
			Mqtt:      line[1],
			Invent:    line[2],
			UnitGUID:  line[3],
			MsgID:     line[4],
			Text:      line[5],
			Context:   line[6],
			Class:     line[7],
			Level:     level,
			Area:      line[9],
			Addr:      line[10],
			Block:     line[11],
			Type:      line[12],
			Bit:       bit,
			InvertBit: invertBit,
		}

		messages = append(messages, currentMessage)
	}

	return messages, err
}

// str to int for columns _n, _level, _bit from file.
// if _bit is empty -> _bit = -1
// if _invertBit is empty -> _invertBit = -1
func transformDataType(_level, _bit, _invertBit string) (level, bit, invertBit int, err error) {

	level, err = strconv.Atoi(_level)
	if err != nil {
		return level, bit, invertBit, err
	}

	if _bit == "" {
		bit = -1
	} else {
		bit, err = strconv.Atoi(_bit)
	}
	if err != nil {
		return level, bit, invertBit, err
	}

	if _invertBit == "" {
		invertBit = -1
	} else {
		invertBit, err = strconv.Atoi(_invertBit)
	}
	if err != nil {
		return level, bit, invertBit, err
	}
	return level, bit, invertBit, err
}

func processFileName(filePath string) (newPath string, err error) {
	directory, fileName := filepath.Split(filePath)
	newFilename := fmt.Sprintf("processed_%s", fileName)
	newPath = filepath.Join(directory, newFilename)

	err = os.Rename(filePath, newPath)
	if err != nil {
		return newPath, err
	}

	return newPath, nil
}

func unProcessFileName(filePath string) (newFilename string, err error) {
	directory, fileName := filepath.Split(filePath)
	originalFilename := strings.TrimPrefix(fileName, "processed_")
	newPath := filepath.Join(directory, originalFilename)

	err = os.Rename(filePath, newPath)
	if err != nil {
		return newPath, err
	}

	return newPath, nil
}
