package main

import "fmt"

type httpError struct {
	statusCode     int
	msg, msgPublic string
}

func newHTTPError(statusCode int, msg, msgPublic string, orgErr error) httpError {
	if orgErr != nil {
		return httpError{statusCode: statusCode, msg: fmt.Sprintf("%s, org error: %s", msg, orgErr.Error()), msgPublic: msgPublic}
	}

	return httpError{statusCode: statusCode, msg: msg, msgPublic: msgPublic}
}

func (he httpError) Error() string {
	return he.msg
}
