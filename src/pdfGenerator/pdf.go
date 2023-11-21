package pdfGenerator

import (
	"github.com/fumiama/go-docx"
	"os"
)

func F(text string) {

	w := docx.NewA4()

	// Добавляем параграф с текстом
	para1 := w.AddParagraph()
	para1.AddText(text).AddTab()
	para1.AddText("size").Size("44").AddTab()
	f, err := os.Create("generated.docx")
	// save to file
	if err != nil {
		panic(err)
	}
	_, err = w.WriteTo(f)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}
