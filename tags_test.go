package retag

import (
	"encoding/json"
	"fmt"
	"testing"
)

type User struct {
	ID    string
	Name  string
	Class []Class
	Age   int
}

type Class struct {
	ClassID   string
	ClassName string
	Scores    map[string]int
}

type ClassTag struct {
	ClassID   string         `json:"class_id" xml:"class_id"`
	ClassName string         `json:"class_name" xml:"class_name"`
	Scores    map[string]int `json:"scores"`
}

type UserTag struct {
	ID    string               `json:"_id" xml:"_id"`
	Name  string               `json:"_name" xml:"_name"`
	Class map[string]*ClassTag `json:"class"`
	Age   int                  `json:"_age"`
}

func TestRetag(t *testing.T) {
	// fts := FieldTag{
	// 	"ID":   `json:"_id"`,
	// 	"Name": `json:"name"`,
	// }

	fts := GetFieldTags(&UserTag{})
	fmt.Println(fts)
	u := User{
		"001",
		"qianlnk",
		[]Class{
			{"01", "math", map[string]int{"math": 100}},
			{"02", "english", map[string]int{"english": 100}},
		},
		18,
	}
	nu := Retag(&u, fts)

	data, err := json.Marshal(nu)
	fmt.Println(string(data), err)
}
