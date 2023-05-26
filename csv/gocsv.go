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
	wr.WriteString(strings.Join(titles, ","))
	wr.WriteString("\n")
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
		titles = append(titles, line.Name)
	}

	return titles
}

func UnmarshalCSV[T any](csv []byte, target *[]T) error {
	for i, line := range bytes.Split(csv, []byte("\n")) {
		splitted := bytes.Split(line, []byte(","))
		if i != 0 {
			for _, v := range splitted {
				fmt.Println(string(v))
				var t T
				val := reflect.TypeOf(target).Kind()
				fmt.Println(val)

				switch val {
				case reflect.String:
					fmt.Println("string")
					reflect.ValueOf(t).Elem().SetString(reflect.ValueOf(v).String())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					reflect.ValueOf(t).Elem().SetInt(reflect.ValueOf(v).Int())
					fmt.Println("int")
				default:
					fmt.Printf("unhandled kind %s", val)
				}
			}
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
