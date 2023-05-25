package main

import (
	"gocsv/csv"
)

type Person struct {
	Name string `csv:"name"`
	Age  int    `csv:"age"`
}

func (p Person) MarshalCSV() ([]byte, error) {
	return []byte(p.Name + "," + string(p.Age)), nil
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
}
