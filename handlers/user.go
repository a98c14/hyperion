package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type UserHandler struct{}


func (h *UserHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	time.Sleep(2 * time.Second)
	fmt.Println(req.URL.Path)
}
