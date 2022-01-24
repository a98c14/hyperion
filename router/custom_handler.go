package router

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/a98c14/hyperion/common/errors"
	"github.com/a98c14/hyperion/common/response"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO(selim): Add response based on type
	if err := h(w, r); err != nil {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		response.InternalError(w, errors.Wrap(name, err))
		return
	}
}
