package dbInterface

import (
	"Biocad2/src/message"
	"github.com/spf13/viper"
)

type DBI interface {
	Init(config *viper.Viper)
	AddFile(messages []message.Message, filePath string) error
	AddFileName(FileName string, state string) error
	GetById(UnitGUID string, pageNumber, pageSize int) ([]message.Message, int64)
	AllGetById(UnitGUID string) []message.Message
}
