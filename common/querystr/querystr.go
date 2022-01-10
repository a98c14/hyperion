package querystr

import "strconv"

type Identifiable interface {
	Id(i int) int
	Len() int
}

// TODO(selim): Use below return values instead of string
func GenerateInString(ints []int) string {
	ids := ""
	for idx, v := range ints {
		if idx != 0 {
			ids += "," + strconv.Itoa(v)
		} else {
			ids += strconv.Itoa(v)
		}
	}
	return ids
}

func GetIntArray(values Identifiable) []int {
	arr := make([]int, values.Len())
	for i := 0; i < values.Len(); i++ {
		arr[i] = values.Id(i)
	}
	return arr
}

// Generates placeholder value string and values to write to in place for pgx
// Example:
// str -> $1, $2, $3
// arr -> 4, 6, 7
func GenerateInStringIdentifiable(values Identifiable) (string, []interface{}) {
	str := ""
	arr := make([]interface{}, values.Len())
	for i := 0; i < values.Len(); i++ {
		v := values.Id(i)
		arr[i] = v
		str += "$" + strconv.Itoa(i+1) + ","
	}
	str = str[:len(str)-1] // Remove last `,`
	return str, arr
}
