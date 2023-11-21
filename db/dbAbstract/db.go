package dbAbstract

import "biocadGo/src/message"

type Database interface {
	Init(host, port string)
	AddFile(messages []message.Message) error
	GetById(UnitGUID string, pageNumber, pageSize int) ([]message.Message, int64)
}
