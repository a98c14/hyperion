package errors

import stderror "errors"

var ErrExists = stderror.New("given id already exists")

var ErrBadRequest = stderror.New("can't parse request body")

var ErrNotFound = stderror.New("resource can not be found")
