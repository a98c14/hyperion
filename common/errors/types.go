package errors

import stderror "errors"

var ExistsError = stderror.New("Given Id already exists")
