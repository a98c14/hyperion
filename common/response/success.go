package response

import (
	"net/http"

	"github.com/a98c14/hyperion/common/json"
)

func Json(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	err := json.Encode(w, &body)
	if err != nil {
		InternalError(w, err)
	}
}
