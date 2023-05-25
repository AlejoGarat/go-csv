package csv

import (
	"bytes"
	"reflect"
)

type MarshalCSVer interface {
	MarshalCSV() ([]byte, error)
}

type UnmarshalCVer interface {
	UnmarshalCSV([]byte) error
}

func MarshalCSV(arr []any) ([]byte, error) {
	// titles := getTitlesWithReflection(arr)
	wr := bytes.Buffer{}
	// wr.WriteString(strings.Join(titles, ","))
	// wr.WriteString("\n")
	for _, v := range arr {
		if m, ok := v.(MarshalCSVer); ok {
			bs, err := m.MarshalCSV()
			if err != nil {
				return nil, err
			}
			wr.Write(bs)
			wr.WriteString("\n")
		} else {
			panic("not implemented")
		}
	}

	return nil, nil
}

func getTitlesWithReflection(p []any) []string {
	v := reflect.ValueOf(p[0])
	y := v.Interface().([]string)
	return y
}

func UnMarshalCSV([]byte, []any) error {
	return nil
}
