package response

import (
	"net/http"

	"github.com/a98c14/hyperion/common/json"
)

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
func Json(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	err := json.Encode(w, &body)
	if err != nil {
		InternalError(w, err)
	}
}

func Success(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	err := json.Encode(w, struct{ Message string }{Message: message})
	if err != nil {
		InternalError(w, err)
	}
}
