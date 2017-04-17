package retag

import (
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

func getElemType(t reflect.Type) reflect.Type {
	for {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		} else {
			return t
		}
	}
}

func GetFieldTags(s interface{}) FieldTag {
	t := reflect.TypeOf(s)

	t = getElemType(t)

	return getTag(t)

}

func addFieldTag(a FieldTag, b FieldTag, prefix string) FieldTag {
	res := a
	for f, t := range b {
		res[FieldName(prefix+string(f))] = t
	}

	return res
}

func getTag(t reflect.Type) FieldTag {
	fts := make(FieldTag)

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			fts[FieldName(t.Field(i).Name)] = t.Field(i).Tag
			subfts := getTag(getElemType(t.Field(i).Type))
			fts = addFieldTag(fts, subfts, t.Field(i).Name+".")
		}
	case reflect.Array, reflect.Slice, reflect.Map:
		subfts := getTag(getElemType(t.Elem()))
		fts = addFieldTag(fts, subfts, "")
	default:
		//fmt.Println(t.String())
	}

	return fts
}
