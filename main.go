package main

import (
	"fmt"
	"strconv"

	"gocsv/csv"
)

type Person struct {
	Name string `csv:"Pepe"`
	Age  int    `csv:"-"`
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

	data := []byte(`Pepe, Name
	hola,John
	cantu no sabe rust,Jane`)
	var unmarshalled []Person

	err = csv.UnmarshalCSV(data, &unmarshalled)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(unmarshalled)
	}
}
