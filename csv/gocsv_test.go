package csv_test

import (
	"strconv"
	"testing"

	"gocsv/csv"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string `csv:"name"`
	Age  int    `csv:"age"`
}

type PersonIgnoredTag struct {
	Name string `csv:"name"`
	Age  int    `csv:"-"`
}

type PersonTagPriority struct {
	Name string `csv:"some_tag"`
	Age  int    `csv:"name"`
}

type PersonChangedOrder struct {
	Age  int    `csv:"age"`
	Name string `csv:"name"`
}

type Student struct {
	Name   string
	Number int
}

type StudentTag struct {
	Name   string `json:"name"`
	Number int    `json:"age"`
}

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

func (p Person) MarshalCSV() (map[string][]byte, error) {
	return map[string][]byte{
		"name": []byte(p.Name),
		"age":  []byte(strconv.Itoa(p.Age)),
	}, nil
}

func (p Person) UnmarshalCSV([]byte) error {
	return nil
}

func getPersons() []Person {
	return []Person{
		{"John", 20},
		{"Jane", 30},
	}
}

func getPersonsIgnoredTag() []PersonIgnoredTag {
	return []PersonIgnoredTag{
		{"John", 20},
		{"Jane", 30},
	}
}

func getStudents() []Student {
	return []Student{
		{"John", 1234},
		{"Jane", 8765},
	}
}

func TestMarshalWithTags(t *testing.T) {
	persons := getPersons()
	b, err := csv.MarshalCSV([]any{persons[0], persons[1]})
	assert.NoError(t, err)
	assert.Equal(t, "name,age\nJohn,20\nJane,30\n", string(b))
}

func TestMarshalWithoutTags(t *testing.T) {
	students := getStudents()
	b, err := csv.MarshalCSV([]any{students[0], students[1]})
	assert.NoError(t, err)
	assert.Equal(t, "Name,Number\nJohn,1234\nJane,8765\n", string(b))
}

func TestMarshalWithoutCsvTags(t *testing.T) {
	student := StudentTag{
		"John", 1234,
	}
	b, err := csv.MarshalCSV([]any{student})
	assert.NoError(t, err)
	assert.Equal(t, "Name,Number\nJohn,1234\n", string(b))
}

func TestMarshalWithIgnoredTag(t *testing.T) {
	persons := getPersonsIgnoredTag()
	b, err := csv.MarshalCSV([]any{persons[0], persons[1]})
	assert.NoError(t, err)
	assert.Equal(t, "name\nJohn\nJane\n", string(b))
}

func TestMarshalWithMixedOrder(t *testing.T) {
	arr := []any{struct{ Name, Age string }{"John", "20"}, struct{ Age, Name string }{"30", "Jane"}}
	b, err := csv.MarshalCSV(arr)
	assert.NoError(t, err)
	assert.Equal(t, "Name,Age\nJohn,20\nJane,30\n", string(b))
}

func TestMarshalStructWithPointer(t *testing.T) {
	type typ struct{ Name *string }
	stringPointer := func(s string) *string { return &s }
	arr := []typ{{Name: stringPointer("John")}}

	b, err := csv.MarshalCSV(arr)
	assert.NoError(t, err)
	assert.Equal(t, "Name\nJohn\n", string(b))
}

func TestUnmarshalWithTags(t *testing.T) {
	data := []byte(`name, age
					John, 20
					Jane, 30`)
	var unmarshalled []Person

	err := csv.UnmarshalCSV(data, &unmarshalled)
	assert.NoError(t, err)
	assert.Equal(t, getPersons(), unmarshalled)
}

func TestUnmarshalWithoutTags(t *testing.T) {
	data := []byte(`Name, Number
					John, 1234
					Jane, 8765`)
	var unmarshalled []Student

	err := csv.UnmarshalCSV(data, &unmarshalled)
	assert.NoError(t, err)
	assert.Equal(t, getStudents(), unmarshalled)
}

func TestUnmarshalWithIgnoredTag(t *testing.T) {
	data := []byte(`name
					John
					Jane`)
	var unmarshalled []PersonIgnoredTag

	err := csv.UnmarshalCSV(data, &unmarshalled)
	assert.NoError(t, err)

	personIgnoredTag := []PersonIgnoredTag{
		{"John", 0},
		{"Jane", 0},
	}

	assert.Equal(t, personIgnoredTag, unmarshalled)
}

func TestUnmarshalTagPriority(t *testing.T) {
	data := []byte(`some_tag, name
					John, 	  20
					Jane, 	  30`)
	var unmarshalled []PersonTagPriority

	err := csv.UnmarshalCSV(data, &unmarshalled)
	assert.NoError(t, err)

	personTagPriority := []PersonTagPriority{
		{"John", 20},
		{"Jane", 30},
	}

	assert.Equal(t, personTagPriority, unmarshalled)
}

func TestBigMarshalUnmarshal(t *testing.T) {
	var s []any
	for i := 0; i < 10; i++ {
		s = append(s, SomeStruct{
			"Alejo", 10, "Martin", 30, "Cantu", 50, "Pepe", 70, "Fede", 90,
		})
	}

	b, err := csv.MarshalCSV(s)
	assert.NoError(t, err)

	var unmarshalled []SomeStruct
	err = csv.UnmarshalCSV(b, &unmarshalled)
	assert.NoError(t, err)

	for i, x := range unmarshalled {
		if i < len(unmarshalled) {
			assert.Equal(t, x.A, "Alejo")
			assert.Zero(t, x.B)
			assert.Equal(t, x.C, "Martin")
			assert.Zero(t, x.D)
			assert.Zero(t, x.E)
			assert.Equal(t, x.F, 50)
			assert.Zero(t, x.G)
			assert.Equal(t, x.H, 70)
			assert.Equal(t, x.I, "Fede")
			assert.Zero(t, x.J)
		}
	}
}
