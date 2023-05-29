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
	}
	for _, v := range arr {
		getRowsWithReflection(v, &wr)
	}

	return wr.Bytes(), nil
}

func getRowsWithReflection(p any, wr *bytes.Buffer) {
	v := reflect.TypeOf(p)
	if v.Kind() != reflect.Struct {
		panic("struct is required")
	} else {
		s := reflect.ValueOf(p)
		for i := 0; i < v.NumField(); i++ {
			line := v.Field(i)
			tag, ok := line.Tag.Lookup("csv")
			if ok && tag != "-" {
				switch f := s.FieldByName(line.Name); f.Kind() {
				case reflect.String:
					wr.WriteString(f.String())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					wr.WriteString(strconv.FormatInt(f.Int(), 10))
				default:
					fmt.Printf("unhandled kind %s", v.Kind())
				}
				if i < v.NumField()-1 {
					wr.WriteString(",")
				}
			}
		}
		wr.WriteString("\n")
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
		} else if tag != "-" {
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
			titles = strings.Split(string(line), ",")
			for i := range titles {
				titles[i] = strings.TrimSpace(titles[i])
			}
			var t T
			typ := reflect.TypeOf(t)
			for i := 0; i < typ.NumField(); i++ {
				line := typ.Field(i)
				fields = append(fields, line)
			}
		} else {
			splitted := bytes.Split(line, []byte(","))
			var t T
			for j, s := range splitted {
				val := reflect.ValueOf(&t)
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
				if index == -1 {
					continue
				}

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
			*target = append(*target, t)
		}
	}
	return nil
}

/*
if m, ok := v.(MarshalCSVer); ok {
			bs, err := m.MarshalCSV()
			if err != nil {
				return nil, err
			}
			wr.Write(bs)
			wr.WriteString("\n")
*/
