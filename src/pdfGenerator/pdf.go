package pdfGenerator

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
	"unicode/utf8"
)

func ObjToStrArr(msg interface{}) (result []string) {
	val := reflect.ValueOf(msg)
	typ := reflect.TypeOf(msg)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		fieldValue := fmt.Sprintf("%v", field.Interface())
		if fieldName != "UnitGUID" {
			result = append(result, fieldValue)
		}
	}
	return result
}

func ObjHeadToStrArr(msg interface{}) (result []string) {
	val := reflect.ValueOf(msg)
	typ := reflect.TypeOf(msg)

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		if fieldName != "UnitGUID" {
			result = append(result, fieldName)
		}
	}
	return result
}

func MakeReportPDF(PDFConfig *viper.Viper, Header string, headers []string, data [][]string) error {

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFooterFunc(func() {

		pdf.SetY(-15) // Position at 1.5 cm from bottom
		pdf.SetFont("Arial", "I", 8)

		pdf.CellFormat(10, 10, fmt.Sprint(pdf.PageNo()), "", 0, "L", false, 0, "") // Page number
	})
	fontData, err := os.ReadFile(PDFConfig.GetString("font_path"))
	if err != nil {
		log.Println("Ошибка при чтении файла шрифта:", err)
		return err
	}
	pdf.AddUTF8FontFromBytes("ArialUnicode", "", fontData)

	pdf.AddPage()
	{
		pdf.SetFont(PDFConfig.GetString("headerFont"), "", PDFConfig.GetFloat64("headerSize"))

		// Define layers
		l1 := pdf.AddLayer("Layer 1", true)
		pdf.BeginLayer(l1)
		pdf.Write(8, Header+"\n")
		pdf.EndLayer()
	}

	// Устанавливаем шрифт и размер текста для заголовка таблицы
	pdf.SetFont(PDFConfig.GetString("tableFont"), "", PDFConfig.GetFloat64("tableHeaderSize"))

	// расчет размера ячейки
	var sizeArr []float64
	{
		var columnLen int
		for i := range headers {
			columnLen = utf8.RuneCountInString(headers[i])
			dataLen := utf8.RuneCountInString(data[0][i])

			if columnLen > dataLen {
				sizeArr = append(sizeArr, float64(columnLen))
			} else {
				sizeArr = append(sizeArr, float64(dataLen))
			}
		}
	}

	// Создаем таблицу
	{
		for i, cell := range headers {
			pdf.CellFormat(sizeArr[i]*2.1, 7, cell, PDFConfig.GetString("borderType"), 0, "1", false, 0, "")
		}
		pdf.Ln(-1)

		// Устанавливаем шрифт и размер текста для тела таблицы
		pdf.SetFont(PDFConfig.GetString("tableFont"), "", PDFConfig.GetFloat64("tableSize"))
		for _, row := range data {
			for i, cell := range row {
				pdf.CellFormat(sizeArr[i]*2.1, 7, cell, PDFConfig.GetString("borderType"), 0, "1", false, 0, "")
			}
			// Переходим на новую строку для следующей строки таблицы
			pdf.Ln(-1)
		}
	}

	// Сохраняем PDF в файл
	err = pdf.OutputFileAndClose(PDFConfig.GetString("report_path") + "/" + Header + ".pdf")
	if err != nil {
		log.Println("Ошибка при сохранении файла:", err)
		return err
	}
	return err
}

func PDF() {}
