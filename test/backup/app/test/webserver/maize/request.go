package maize

import (
	"./session"
	"net/http"
)

type Request struct {
	*http.Request
	Session *session.SessionStore
}
