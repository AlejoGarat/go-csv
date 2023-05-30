package csv

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type MarshalCSVer interface {
	MarshalCSV() (map[string][]byte, error)
}

type UnmarshalCVer interface {
	UnmarshalCSV([]byte) error
}

func MarshalCSV(arr []any) ([]byte, error) {
	titles := getTitlesWithReflection(arr)
	wr := bytes.Buffer{}
	if len(titles) != 0 {
		wr.WriteString(strings.Join(titles, ","))
		wr.WriteString("\n")
		getRowsWithReflection(arr, &wr, titles)
	}

	return wr.Bytes(), nil
}

func getRowsWithReflection(p []any, wr *bytes.Buffer, titles []string) {
	v := reflect.TypeOf(p[0])
	if v.Kind() != reflect.Struct {
		panic("struct is required")
	} else {
		for _, s := range p {
			fields := getFields(v)

			v := reflect.TypeOf(s)
			val := reflect.ValueOf(s)

			for i, field := range fields {
				if field.Tag.Get("csv") != "-" {
					if i > 0 {
						wr.WriteString(",")
					}
					switch f := val.FieldByName(field.Name); f.Kind() {
					case reflect.String:
						wr.WriteString(f.String())
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						wr.WriteString(strconv.FormatInt(f.Int(), 10))
					default:
						fmt.Printf("unhandled kind %s", v.Kind())
					}
				}
			}

			wr.WriteString("\n")
		}
	}
}

func getTitlesWithReflection(p []any) []string {
	var titles []string

	v := reflect.TypeOf(p[0])
	if v.Kind() != reflect.Struct {
		panic("struct is required")
	}

	for i := 0; i < v.NumField(); i++ {
		line := v.Field(i)

		tag, ok := line.Tag.Lookup("csv")
		if ok && tag != "-" {
			titles = append(titles, tag)
		} else if !ok {
			titles = append(titles, line.Name)
		}
	}

	return titles
}

func UnmarshalCSV[T any](csv []byte, target *[]T) error {
	var titles []string
	var fields []reflect.StructField

	for i, line := range bytes.Split(csv, []byte("\n")) {
		if i == 0 {
			titles = getTitles(line)
			var t T
			typ := reflect.TypeOf(t)
			fields = getFields(typ)
			continue
		}
		splitted := bytes.Split(line, []byte(","))
		if len(splitted) == 0 || len(splitted) == 1 && len(bytes.TrimSpace(splitted[0])) == 0 {
			continue
		}
		var t T
		for j, s := range splitted {
			val := reflect.ValueOf(&t)

			index := updateIndex(fields, titles, j)
			if index == -1 {
				continue
			}

			setField(index, val, s)

		}
		*target = append(*target, t)
	}
	return nil
}

func getTitles(line []byte) []string {
	titles := strings.Split(string(line), ",")
	for i := range titles {
		titles[i] = strings.TrimSpace(titles[i])
	}
	return titles
}

func getFields(typ reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	for i := 0; i < typ.NumField(); i++ {
		line := typ.Field(i)
		fields = append(fields, line)
	}

	return fields
}

func updateIndex(fields []reflect.StructField, titles []string, j int) int {
	index := -1

	for i, field := range fields {
		if field.Tag.Get("csv") == titles[j] {
			index = i
			break
		}
	}

	if index == -1 {
		for i, field := range fields {
			if field.Name == titles[j] {
				if _, ok := field.Tag.Lookup("csv"); !ok {
					index = i
					break
				}
			}
		}
	}

	return index
}

func setField(index int, val reflect.Value, s []byte) {
	fieldValue := val.Elem().Field(index)
	k := fieldValue.Kind()

	switch k {
	case reflect.String:
		val.Elem().Field(index).SetString(strings.TrimSpace(string(s)))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		is, err := strconv.ParseUint(strings.TrimSpace(string(s)), 10, 64)
		if err != nil {
			panic(err.Error())
		}

		val.Elem().Field(index).SetInt(int64(is))

	default:
		fmt.Printf("unhandled kind %s", val)
	}
}
