package upnp

import (
// "log"
)

type Msg struct {
	request  Request
	response Response
}

type Request interface {
	send()
}

type Response interface {
	resolve()
}
