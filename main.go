package main

import (
	"strconv"

	"gocsv/csv"
)

type Person struct {
	Name string `csv:"name"`
	Age  int    `csv:"age"`
}

func (p Person) MarshalCSV() (map[string][]byte, error) {
	return map[string][]byte{
		"name": []byte(p.Name),
		"age":  []byte(strconv.Itoa(p.Age)),
	}, nil
}

func (p Person) UnmarshalCSV([]byte) error {
	return nil
}

func main() {
	p1 := Person{"John", 20}
	p2 := Person{"Jane", 30}

	b, err := csv.MarshalCSV([]any{p1, p2})
	if err != nil {
		panic(err)
	}

	println(string(b))

	data := []byte(`Age, Name
	20, John
	30, Jane`)
	var unmarshalled []Person

	err = csv.UnmarshalCSV(data, &unmarshalled)
	if err != nil {
		println(err.Error())
	} else {
		println(unmarshalled)
	}
}
