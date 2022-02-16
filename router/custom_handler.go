package router

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/a98c14/hyperion/common"
	"github.com/a98c14/hyperion/common/errors"
	"github.com/a98c14/hyperion/common/response"
)

type Handler func(state common.State, w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state, err := common.InitState(r)
	if err != nil {
		response.ErrorWhileInitializing(w, err)
		return
	}

	if err := h(state, w, r); err != nil {
		switch err {
		case errors.ErrExists:
			response.BadRequest(w, err)
		default:
			name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
			response.InternalError(w, errors.Wrap(name, err))
			return
		}
	}
}
