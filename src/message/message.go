package message

import (
	"fmt"
	"reflect"
	"strings"
)

type Message struct {
	Mqtt      string `json:"mqtt"`
	Invent    string `json:"invent"`
	UnitGUID  string `json:"unit_guid"`
	MsgID     string `json:"msg_id"`
	Text      string `json:"text"`
	Context   string `json:"context"`
	Class     string `json:"class"`
	Level     int    `json:"level"`
	Area      string `json:"area"`
	Addr      string `json:"addr"`
	Block     string `json:"block"`
	Type      string `json:"type"`
	Bit       int    `json:"bit"`
	InvertBit int    `json:"invert_bit"`
}

func (m *Message) ToString() string {
	v := reflect.ValueOf(m)

	var fields []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		value := fmt.Sprintf("%v", field.Interface())
		fields = append(fields, fmt.Sprintf("%s", value))
	}
	return strings.Join(fields, ", ")
}
