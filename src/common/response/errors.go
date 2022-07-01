package response

import (
	"fmt"
	"net/http"

	"github.com/a98c14/hyperion/common/json"
)

type resp struct {
	Message string
}

func InternalError(w http.ResponseWriter, err error) {
	fmt.Println("Internal Error: " + err.Error())
	wrapped := resp{
		Message: err.Error(),
	}
	w.WriteHeader(http.StatusInternalServerError)
	json.Encode(w, &wrapped)
}

func BadRequest(w http.ResponseWriter, err error) {
	fmt.Println("Bad Request: ", err.Error())
	wrapped := resp{
		Message: err.Error(),
	}
	w.WriteHeader(http.StatusBadRequest)
	json.Encode(w, &wrapped)
}

func ErrorWhileInitializing(w http.ResponseWriter, err error) {
	fmt.Println("Error while initializing: " + err.Error())
	http.Error(w, "Error while initializing. "+err.Error(), http.StatusInternalServerError)
}
