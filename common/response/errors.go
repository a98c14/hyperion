package response

import (
	"fmt"
	"net/http"
)

func InternalError(w http.ResponseWriter, err error) {
	fmt.Println("Internal Error: " + err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, err error) {
	fmt.Println("Bad Request: ", err.Error())
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func ErrorWhileInitializing(w http.ResponseWriter, err error) {
	fmt.Println("Error while initializing: " + err.Error())
	http.Error(w, "Error while initializing. "+err.Error(), http.StatusInternalServerError)
}
