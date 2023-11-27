package pdfGenerator

import (
	"reflect"
	"testing"
)

func TestObjToStrArr(t *testing.T) {
	type MyStruct struct {
		StringField     string
		IntField        int
		FloatField      float64
		BoolField       bool
		UintField       uint
		Int8Field       int8
		Uint8Field      uint8
		Int16Field      int16
		Uint16Field     uint16
		Int32Field      int32
		Uint32Field     uint32
		Int64Field      int64
		Uint64Field     uint64
		ByteField       byte
		RuneField       rune
		Float32Field    float32
		Complex64Field  complex64
		Complex128Field complex128
	}

	type testCase struct {
		name       string
		args       interface{}
		wantResult []string
	}
	instance := MyStruct{
		StringField:     "Hello",
		IntField:        42,
		FloatField:      3.14,
		BoolField:       true,
		UintField:       100,
		Int8Field:       -10,
		Uint8Field:      20,
		Int16Field:      -1000,
		Uint16Field:     2000,
		Int32Field:      -100000,
		Uint32Field:     200000,
		Int64Field:      -100000000,
		Uint64Field:     200000000,
		ByteField:       'x',
		RuneField:       'Ω',
		Float32Field:    1.234,
		Complex64Field:  complex(1, 2),
		Complex128Field: complex(3, 4),
		// Установите значения для других полей по аналогии с их типами
	}
	testCases := []testCase{
		{
			name:       "all base type GO",
			args:       instance,
			wantResult: []string{"Hello", "42", "3.14", "true", "100", "-10", "20", "-1000", "2000", "-100000", "200000", "-100000000", "200000000", "120", "937", "1.234", "(1+2i)", "(3+4i)"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ObjToStrArr(tt.args); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ObjToStrArr() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

// Тест для функции ObjHeadToStrArr
func TestObjHeadToStrArr(t *testing.T) {
	type TestStruct struct {
		Field1   string
		Field2   int
		Field3   bool
		UnitGUID string
		Field4   float64
	}

	testCases := []struct {
		input    interface{}
		expected []string
	}{
		{
			input: TestStruct{},
			expected: []string{
				"Field1",
				"Field2",
				"Field3",
				"Field4",
			},
		},
		// Добавьте больше тестовых случаев, если необходимо
	}

	for _, testCase := range testCases {
		result := ObjHeadToStrArr(testCase.input)
		if len(result) != len(testCase.expected) {
			t.Errorf("Неверная длина результата. Ожидалось %d, получено %d", len(testCase.expected), len(result))
		}
		for i, val := range result {
			if val != testCase.expected[i] {
				t.Errorf("Неверное значение на позиции %d. Ожидалось %s, получено %s", i, testCase.expected[i], val)
			}
		}
	}
}
