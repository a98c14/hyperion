package json

import (
	j "encoding/json"
	"io"
	"net/http"
)

func Decode(r *http.Request, v interface{}) error {
	return j.NewDecoder(r.Body).Decode(&v)
}

func Encode(w io.Writer, v interface{}) error {
	return j.NewEncoder(w).Encode(&v)
}
