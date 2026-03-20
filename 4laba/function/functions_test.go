package function

import (
	"auiapp/model"
	"reflect"
	"testing"
	"time"
)

func TestParcingLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected model.ProjectWork
	}{
		{
			name:  "Valid line",
			input: `"Иван Иванов" "тесты" 2023.10.25 a`,
			expected: model.ProjectWork{
				Name:       "Иван Иванов",
				NameOfWork: "тесты",
				Date:       time.Date(2023, 10, 25, 0, 0, 0, 0, time.UTC),
				Type:       "a",
			},
		},
		{
			name:     "Invalid format",
			input:    `Иван "Работа" 2023.10.25 b`,
			expected: model.ProjectWork{},
		},
		{
			name:     "Empty",
			input:    "",
			expected: model.ProjectWork{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := ParcingLine(test.input)
			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("ParcingLine(%q) => %v, want %v", test.input, res, test.expected)
			}
		})
	}
}

func TestParcingFile(t *testing.T) {
	data := []byte(`"Иван Иванов" "UnitTest1" 2026.12.12 a
"Иван Иванов 2" "Test2" 2026.12.12 f`)
	results := ParsingFile(data)
	if len(results) != 2 {
		t.Errorf("Need 2 string but have %d ", len(results))
	}
}
