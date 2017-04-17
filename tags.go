package retag

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type FieldName string
type FieldTag map[FieldName]reflect.StructTag

func Retag(s interface{}, fts FieldTag) interface{} {
	ptrVal := reflect.ValueOf(s)

	newType := getType(ptrVal.Elem().Type(), fts)
	newPtrVal := reflect.NewAt(newType, unsafe.Pointer(ptrVal.Pointer()))

	return newPtrVal.Interface()
}

func getType(t reflect.Type, fts FieldTag) reflect.Type {
	switch t.Kind() {
	case reflect.Struct:
		return getStructType(t, fts)
	case reflect.Ptr:
		return reflect.PtrTo(getType(t.Elem(), fts))
	case reflect.Array:
		return reflect.ArrayOf(t.Len(), getType(t.Elem(), fts))
	case reflect.Slice:
		return reflect.SliceOf(getType(t.Elem(), fts))
	case reflect.Map:
		return reflect.MapOf(getType(t.Key(), fts), getType(t.Elem(), fts))
	case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Interface:
		panic("Unsupported type: " + t.Kind().String())
	default:
		return t
	}
}

func getStructType(t reflect.Type, fts FieldTag) reflect.Type {
	if t.NumField() == 0 || len(fts) == 0 {
		return t
	}

	fields := make([]reflect.StructField, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if isExported(field.Name) {
			subfts := make(FieldTag)
			for f, t := range fts {
				if strings.HasPrefix(string(f), field.Name) {
					subfts[FieldName(strings.TrimPrefix(string(f), field.Name+"."))] = t
				}
			}

			field.Type = getType(field.Type, subfts)
		} else {
			//BUG with unexport field
			field.Name = ""
			field.PkgPath = ""
		}
		if _, ok := fts[FieldName(field.Name)]; ok {
			field.Tag = fts[FieldName(field.Name)]
		}

		fields = append(fields, field)
	}

	newType := reflect.StructOf(fields)
	//fmt.Println(newType)
	if t.Size() != newType.Size() {
		panic("Unexpected case, new type has a size different from size of original type")
	}

	return newType
}

func isExported(name string) bool {
	b := name[0]
	return !('a' <= b && b <= 'z') && b != '_'
}

func GetFieldTags(s interface{}) FieldTag {
	ptrVal := reflect.ValueOf(s)

	return getTag(ptrVal.Type().Elem())

}

func getTag(t reflect.Type) FieldTag {
	if t.NumField() == 0 {
		return nil
	}

	fts := make(FieldTag)

	for i := 0; i < t.NumField(); i++ {
		subfts := getStructTag(t.Field(i))
		for f, t := range subfts {
			fts[f] = t
		}
	}

	return fts
}

func getStructTag(f reflect.StructField) FieldTag {
	fts := make(FieldTag)

	fmt.Println(f.Type.Kind())

	switch getType(f.Type, nil).Kind() {
	case reflect.Struct:
		fts[FieldName(f.Name)] = f.Tag
		for i := 0; i < f.Type.Elem().NumField(); i++ {
			subfts := getStructTag(f.Type.Elem().Field(i))
			for subf, subt := range subfts {
				fts[FieldName(f.Name+"."+string(subf))] = subt
			}
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		fts[FieldName(f.Name)] = f.Tag
		fmt.Println("&&&", reflect.MapOf(f.Type.Key(), f.Type.Elem()), getType(f.Type.Elem(), nil).Kind())
		if getType(f.Type, nil).Kind() == reflect.Struct {

			for i := 0; i < f.Type.Elem().NumField(); i++ {
				subfts := getStructTag(f.Type.Elem().Field(i))
				for subf, subt := range subfts {
					fts[FieldName(f.Name+"."+string(subf))] = subt
				}
			}
		}
	//case reflect.Slice:
	default:
		fts[FieldName(f.Name)] = f.Tag
	}
	return fts
}
