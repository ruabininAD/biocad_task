package db

import "biocadGo/src/message"

type Database interface {
	Init(host, port string)
	AddFile(messages []message.Message) error
}
