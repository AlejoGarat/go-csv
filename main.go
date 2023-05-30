package main

import (
	"fmt"
	"strconv"

	"gocsv/csv"
)

type SomeStruct struct {
	A string `csv:"a"`
	B int    `csv:"-"`
	C string `csv:"c"`
	D int    `csv:"-"`
	E string `csv:"-"`
	F int    `csv:"f"`
	G string `csv:"-"`
	H int    `csv:"h"`
	I string `csv:"i"`
	J int    `csv:"-"`
}

func (s SomeStruct) MarshalCSV() (map[string][]byte, error) {
	return map[string][]byte{
		"a": []byte(s.A),
		"c": []byte(s.C),
		"f": []byte(strconv.Itoa(s.F)),
		"h": []byte(strconv.Itoa(s.H)),
		"i": []byte(s.I),
	}, nil
}

func (s SomeStruct) UnmarshalCSV([]byte) error {
	return nil
}

func main() {
	var s []any
	for i := 0; i < 10; i++ {
		s = append(s, SomeStruct{
			"John", 10, "Jane", 30, "Jerry", 50, "Tom", 70, "Jessie", 90,
		})
	}

	b, err := csv.MarshalCSV(s)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(b))

	var unmarshalled []SomeStruct

	err = csv.UnmarshalCSV(b, &unmarshalled)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(unmarshalled)
	}
}
