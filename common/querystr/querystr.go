package querystr

import "strconv"

type Identifiable interface {
	Id(i int) int
	Len() int
}

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

func GenerateInStringIdentifiable(values Identifiable) string {
	str := ""
	for i := 0; i < values.Len(); i++ {
		v := values.Id(i)
		if i != 0 {
			str += "," + strconv.Itoa(v)
		} else {
			str += strconv.Itoa(v)
		}
	}
	return str
}
