package retag

import (
	"encoding/json"
	"fmt"
	"testing"
)

type StdType struct {
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
	String  string
	Bool    bool
	Byte    byte
	//Complex64  complex64
	//Complex128 complex128
	Interface interface{}
}

type StdTypeTag struct {
	Int     int     `json:"_int"`
	Int8    int8    `json:"_int8"`
	Int16   int16   `json:"_int16"`
	Int32   int32   `json:"_int32"`
	Int64   int64   `json:"_int64"`
	Uint    uint    `json:"_uint"`
	Uint8   uint8   `json:"_uint8"`
	Uint16  uint16  `json:"_uint16"`
	Uint32  uint32  `json:"_uint32"`
	Uint64  uint64  `json:"_uint64"`
	Float32 float32 `json:"_float32"`
	Float64 float64 `json:"_float64"`
	String  string  `json:"_string"`
	Bool    bool    `json:"_bool"`
	Byte    byte    `json:"_byte"`
	//Complex64  complex64  `json:"_complex64"`
	//Complex128 complex128 `json:"_complex128"`
	Interface interface{} `json:"_interface"`
}

type StdPtrType struct {
	Int       *int
	Uint      *uint
	Float32   *float32
	String    *string
	Bool      *bool
	Byte      *byte
	Complex64 *complex64
	Interface *interface{}
}

type StdPtrTypeTag struct {
	Int       *int         `json:"ptr_int"`
	Uint      *uint        `json:"ptr_uint"`
	Float32   *float32     `json:"ptr_float32"`
	String    *string      `json:"ptr_string"`
	Bool      *bool        `json:"ptr_bool"`
	Byte      *byte        `json:"ptr_byte"`
	Complex64 *complex64   `json:"ptr_complex64"`
	Interface *interface{} `json:"ptr_interface"`
}

type Recursion struct {
	Std      StdType
	PtrStd   *StdType
	ArrayStd [3]StdType
	SliceStd []StdType
	MapStd   map[string]StdType
}

type RecursionTag struct {
	Std      StdTypeTag            `json:"std"`
	PtrStd   *StdTypeTag           `json:"ptr_std"`
	ArrayStd [3]StdTypeTag         `json:"array_std"`
	SliceStd []StdTypeTag          `json:"slice_std"`
	MapStd   map[string]StdTypeTag `json:"map_std"`
}

var std = &StdType{
	Int:     1,
	Int8:    2,
	Int16:   3,
	Int32:   4,
	Int64:   5,
	Uint:    6,
	Uint8:   7,
	Uint16:  8,
	Uint32:  9,
	Uint64:  10,
	Float32: 11.1,
	Float64: 12.2,
	String:  "13",
	Bool:    true,
	Byte:    'a',
	//Complex64:  14.4 + 14i,
	//Complex128: 15.5 + 15i,
	Interface: "test",
}

func TestStdType(t *testing.T) {
	fts := GetFieldTags(&StdTypeTag{})
	fmt.Println(fts)
	stdTag := Retag(std, fts)
	fmt.Println(stdTag)
	data, err := json.MarshalIndent(stdTag, "", "    ")
	fmt.Println(string(data), err)
}

func TestStdPtrType(t *testing.T) {
	var i = 1
	ptrStd := &StdPtrType{
		Int: &i,
	}

	fts := GetFieldTags(&StdPtrTypeTag{})
	stdPtrTag := Retag(ptrStd, fts)
	data, err := json.MarshalIndent(stdPtrTag, "", "    ")
	fmt.Println(string(data), err)
}

func TestRecursion(t *testing.T) {
	rec := &Recursion{
		Std:      *std,
		PtrStd:   std,
		ArrayStd: [3]StdType{*std, StdType{}, StdType{}},
		SliceStd: []StdType{*std},
		MapStd: map[string]StdType{
			"map_std": *std,
		},
	}

	fts := GetFieldTags(RecursionTag{})
	fmt.Println(fts)
	recTag := Retag(rec, fts)
	data, err := json.MarshalIndent(recTag, "", "    ")
	fmt.Println(string(data), err)
}

func TestCustomFieldTag(t *testing.T) {
	fts := FieldTag{
		"Std":      `json:"std"`,
		"Std.Int":  `json:"_int"`,
		"Std.Uint": `json:"_uint"`,
	}

	rec := &Recursion{
		Std:      *std,
		PtrStd:   std,
		ArrayStd: [3]StdType{*std, StdType{}, StdType{}},
		SliceStd: []StdType{*std},
		MapStd: map[string]StdType{
			"map_std": *std,
		},
	}
	fmt.Println(fts)
	recTag := Retag(rec, fts)
	data, err := json.MarshalIndent(recTag, "", "    ")
	fmt.Println(string(data), err)
}
